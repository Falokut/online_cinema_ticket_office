package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/repository"
	account_service "github.com/Falokut/online_cinema_ticket_office/account_service/pkg/account_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/jwt"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/metrics"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AccountService) mustEmbedUnimplementedAccountServiceServer() {}

type AccountService struct {
	account_service.UnimplementedAccountServiceV1Server
	repo                   repository.AccountRepository
	redisRepo              repository.CacheRepo
	logger                 logging.Logger
	nonActivatedAccountTTL time.Duration
	emailWriter            *kafka.Writer
	cfg                    *config.Config
	metrics                metrics.Metrics
}

func (s *AccountService) ShutDown() {
	s.logger.Println("Service shutting down")
	s.repo.ShutDown()
	s.redisRepo.ShutDown()
}
func NewAccountService(repo repository.AccountRepository, logger logging.Logger,
	redisRepo repository.CacheRepo, emailWriter *kafka.Writer,
	cfg *config.Config, metrics metrics.Metrics) *AccountService {
	return &AccountService{repo: repo,
		logger:                 logger,
		redisRepo:              redisRepo,
		nonActivatedAccountTTL: time.Hour,
		emailWriter:            emailWriter,
		cfg:                    cfg,
		metrics:                metrics,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context,
	in *account_service.CreateAccountRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.CreateAccount")
	defer span.Finish()

	if err := validateSignupInput(in); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)
	}
	exist, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if exist {
		return nil, s.createErrorResponce("A user with this email address already exists. "+
			"Please try another one or simple log in.",
			codes.AlreadyExists)
	}

	inCache, err := s.redisRepo.RegistrationCache.IsAccountInCache(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if inCache {
		return nil, s.createErrorResponce("A user with this email address already exists. "+
			"Please try another one or verify email and log in.",
			codes.AlreadyExists)
	}

	err = s.redisRepo.RegistrationCache.CacheAccount(ctx, in.Email,
		repository.CachedAccount{Username: in.Username, Password: in.Password}, s.nonActivatedAccountTTL)

	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	return &emptypb.Empty{}, nil
}

type emailData struct {
	URL      string  `json:"url"`
	MailType string  `json:"mail_type"`
	LinkTTL  float64 `json:"link_TTL"`
}

func (s *AccountService) RequestAccountVerificationToken(ctx context.Context,
	in *account_service.VerificationTokenRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"AccountService.RequestAccountVerificationToken")
	defer span.Finish()

	inAccountDB, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if inAccountDB {
		return nil, s.createErrorResponce("Account with this email is already activated.", codes.AlreadyExists)
	}

	if err := validateEmail(in.Email); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)

	}

	inCache, err := s.redisRepo.RegistrationCache.IsAccountInCache(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if !inCache {
		const status codes.Code = codes.Internal
		s.metrics.IncCacheMiss(int(status), "RequestAccountVerificationToken")
		return nil, s.createErrorResponce("A user with this email address not exist.",
			status)
	}
	s.metrics.IncCacheHits(int(codes.OK), "RequestAccountVerificationToken")

	cfg := config.GetConfig()
	token, err := jwt.GenerateToken(in.Email, cfg.JWT.VerifyAccountToken.Secret, cfg.JWT.VerifyAccountToken.TTL)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	URL := fmt.Sprintf("%s/%s", in.URL, token)
	LinkTTL := cfg.JWT.VerifyAccountToken.TTL.Seconds()
	body, err := json.Marshal(emailData{URL: URL, MailType: "account/activation", LinkTTL: LinkTTL})
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	go func() {
		err = s.emailWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(in.Email),
			Value: body,
		})
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return &emptypb.Empty{}, nil
}

func (s *AccountService) VerifyAccount(ctx context.Context,
	in *account_service.VerifyAccountRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.VerifyAccount")
	defer span.Finish()

	s.logger.Debug("Parsing token")
	email, err := jwt.ParseToken(in.VerificationToken, config.GetConfig().JWT.VerifyAccountToken.Secret)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)
	}

	s.logger.Debug("Checking account existing in cache")
	acc, err := s.redisRepo.RegistrationCache.GetCachedAccount(ctx, email)
	if err != nil {
		const status = codes.NotFound
		s.metrics.IncCacheMiss(int(status), "VerifyAccount")
		return nil, s.createErrorResponce(err.Error(), status)
	}
	s.metrics.IncCacheHits(int(codes.OK), "VerifyAccount")

	s.logger.Debug("Generating hash from password")
	password_hash, err := bcrypt.GenerateFromPassword([]byte(acc.Password), config.GetConfig().Crypto.BcryptCost)
	if err != nil {
		return nil, s.createErrorResponce("Can't generate hash.", codes.Internal)
	}

	account := model.CreateAccountAndProfile{
		Email:            email,
		Username:         acc.Username,
		Password:         string(password_hash),
		RegistrationDate: time.Now(),
	}

	s.logger.Debug("Creating accound and profile")
	if err := s.repo.CreateAccountAndProfile(ctx, account); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	//The error is not critical, the data will still be deleted from the cache.
	if err := s.redisRepo.RegistrationCache.DeleteAccountFromCache(ctx, email); err != nil {
		s.logger.Warning("Can't delete account from registration cache: ", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *AccountService) SignIn(ctx context.Context,
	in *account_service.SignInRequest) (*account_service.AccessResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.SignIn")
	defer span.Finish()

	s.logger.Debug("Getting user by email")
	u, err := s.repo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.NotFound)
	}

	s.logger.Debug("Password and hash comparison")
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)
	}

	s.logger.Debug("Caching session")
	SessionID := uuid.NewString()
	if err := s.redisRepo.SessionsCache.CacheSession(ctx, model.SessionCache{SessionID: SessionID, AccountID: u.UUID, ClientIP: in.ClientIP, LastActivity: time.Now()}); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	return &account_service.AccessResponce{SessionID: SessionID}, nil
}

func (s *AccountService) GetAccountID(ctx context.Context,
	in *emptypb.Empty) (*account_service.AccountID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.GetAccountID")
	defer span.Finish()

	s.logger.Debug("Getting session id from ctx")
	SessionID, err := s.getSessionIDFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, err
	}

	ClientIP, err := s.getClientIPFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getClientIPFromIP: %v", err)
		return nil, err
	}

	s.logger.Debug("Checking session")
	cache, err := s.checkSession(ctx, SessionID, ClientIP)
	if err != nil {
		return nil, err
	}

	go func() {
		s.logger.Debug("Updating last activity for given session")
		if err := s.redisRepo.SessionsCache.UpdateLastActivityForSession(ctx, cache, SessionID, time.Now()); err != nil {
			s.logger.Warning("Session last activity not updated, error: ", err.Error())
		}
	}()

	return &account_service.AccountID{AccountID: cache.AccountID}, nil
}

func (s *AccountService) Logout(ctx context.Context,
	in *emptypb.Empty) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.Logout")
	defer span.Finish()

	s.logger.Debug("Getting session id from ctx")
	SessionID, err := s.getSessionIDFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, err
	}

	ClientIP, err := s.getClientIPFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getClientIPFromIP: %v", err)
		return nil, err
	}

	s.logger.Debug("Checking session")
	cache, err := s.checkSession(ctx, SessionID, ClientIP)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Terminate current session: ", SessionID)
	err = s.redisRepo.SessionsCache.TerminateSessions(ctx, []string{SessionID}, cache.AccountID)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) RequestChangePasswordToken(ctx context.Context,
	in *account_service.ChangePasswordTokenRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.RequestChangePasswordToken")
	defer span.Finish()

	exist, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if !exist {
		return nil, s.createErrorResponce(grpc_errors.ErrNotFound.Error(), codes.NotFound)
	}

	token, err := jwt.GenerateToken(in.Email, s.cfg.JWT.ChangePasswordToken.Secret, s.cfg.JWT.ChangePasswordToken.TTL)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	URL := fmt.Sprintf("%s/%s", in.URL, token)
	LinkTTL := s.cfg.JWT.ChangePasswordToken.TTL.Seconds()
	body, err := json.Marshal(emailData{URL: URL, MailType: "account/forget-password", LinkTTL: LinkTTL})
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	go func() {
		err := s.emailWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(in.Email),
			Value: body,
		})
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return &emptypb.Empty{}, nil
}

func (s *AccountService) ChangePassword(ctx context.Context,
	in *account_service.ChangePasswordRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.ChangePassword")
	defer span.Finish()

	s.logger.Debug("Parsing jwt token")
	email, err := jwt.ParseToken(in.ChangePasswordToken, config.GetConfig().JWT.ChangePasswordToken.Secret)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)
	}

	s.logger.Debug("Checking account existing in DB")
	exist, err := s.repo.IsAccountWithEmailExist(ctx, email)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}
	if !exist {
		return nil, s.createErrorResponce(err.Error(), codes.NotFound)
	}

	s.logger.Debug("Validating incoming password")
	if err := validatePassword(in.NewPassword); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.InvalidArgument)
	}

	s.logger.Debug("Generating hash for incoming password")
	password_hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), config.GetConfig().Crypto.BcryptCost)
	if err != nil {
		return nil, s.createErrorResponce("Can't generate hash.", codes.Internal)
	}

	s.logger.Debug("Changing account password")
	err = s.repo.ChangePassword(ctx, email, string(password_hash))
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) GetAllSessions(ctx context.Context,
	in *emptypb.Empty) (*account_service.AllSessionsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.GetAllSessions")
	defer span.Finish()

	s.logger.Debug("Getting session id from ctx")
	SessionID, err := s.getSessionIDFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, err
	}

	ClientIP, err := s.getClientIPFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getClientIPFromIP: %v", err)
		return nil, err
	}

	s.logger.Debug("Checking session")
	cache, err := s.checkSession(ctx, SessionID, ClientIP)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Getting sessions for account ", cache.AccountID)
	sessions, err := s.redisRepo.SessionsCache.GetSessionsForAccount(ctx, cache.AccountID)
	if err != nil && err != redis.Nil {
		const status = codes.NotFound
		s.metrics.IncCacheMiss(int(status), "GetAllSessions")
		return nil, s.createErrorResponce(err.Error(), status)
	}
	s.metrics.IncCacheHits(int(codes.OK), "GetAllSessions")

	s.logger.Debug("Converting cache data into responce")
	sessionsInfo := make(map[string]*account_service.SessionInfo, len(sessions))
	for key, session := range sessions {
		sessionsInfo[key] = &account_service.SessionInfo{
			ClientIP:     session.ClientIP,
			SessionInfo:  session.SessionInfo,
			LastActivity: timestamppb.New(session.LastActivity),
		}
	}
	s.logger.Debug("Cache data successfully converted into responce, sending responce")

	return &account_service.AllSessionsResponce{Sessions: sessionsInfo}, nil
}

func (s *AccountService) TerminateSessions(ctx context.Context,
	in *account_service.TerminateSessionsRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.TerminateSessions")
	defer span.Finish()

	s.logger.Debug("Getting session id from ctx")
	SessionID, err := s.getSessionIDFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getSessionIDFromCtx: %v", err)
		return nil, err
	}

	ClientIP, err := s.getClientIPFromCtx(ctx)
	if err != nil {
		s.logger.Errorf("getClientIPFromIP: %v", err)
		return nil, err
	}

	s.logger.Debug("Checking session")
	cache, err := s.checkSession(ctx, SessionID, ClientIP)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Terminating sessions")
	if err := s.redisRepo.SessionsCache.TerminateSessions(ctx, in.SessionsToTerminate, cache.AccountID); err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.Internal)
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) createErrorResponce(errorMessage string, statusCode codes.Code) error {
	s.logger.Error(errorMessage)
	return status.Error(statusCode, errorMessage)
}

func (s *AccountService) checkSession(ctx context.Context,
	SessionID string, ClientIP string) (model.SessionCache, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.checkSession")
	defer span.Finish()

	cache, err := s.redisRepo.SessionsCache.GetSessionCache(ctx, SessionID)
	if err != nil {
		const status = codes.NotFound
		s.metrics.IncCacheMiss(int(status), "checkSession")
		return model.SessionCache{}, s.createErrorResponce(err.Error(), status)
	}
	s.metrics.IncCacheHits(int(codes.OK), "checkSession")

	if ClientIP != cache.ClientIP {
		return model.SessionCache{}, s.createErrorResponce("Access denied", codes.PermissionDenied)
	}
	return cache, nil
}

func (s *AccountService) getSessionIDFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext: %v", grpc_errors.ErrNoCtxMetaData)
	}

	sessionID := md.Get("session_id")
	if sessionID[0] == "" {
		return "", status.Errorf(codes.InvalidArgument, "md.Get sessionId: %v", grpc_errors.ErrInvalidSessionId)
	}

	return sessionID[0], nil
}

func (s *AccountService) getClientIPFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "metadata.FromIncomingContext: %v", grpc_errors.ErrNoCtxMetaData)
	}

	sessionID := md.Get("client_ip")
	if sessionID[0] == "" || len(sessionID) == 0 {
		return "", status.Errorf(codes.InvalidArgument, "md.Get sessionId: %v", grpc_errors.ErrInvalidClientIP)
	}

	return sessionID[0], nil
}
