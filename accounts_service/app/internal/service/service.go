package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/repository"
	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/jwt"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/metrics"
	profiles_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/profiles_service/v1/protos"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountService struct {
	accounts_service.UnimplementedAccountsServiceV1Server
	repo                   repository.AccountRepository
	redisRepo              repository.CacheRepo
	logger                 *logrus.Logger
	nonActivatedAccountTTL time.Duration
	emailWriter            *kafka.Writer
	cfg                    *config.Config
	metrics                metrics.Metrics
	errorHandler           errorHandler
	profilesService        profiles_service.ProfilesServiceV1Client
}

func NewAccountService(repo repository.AccountRepository, logger *logrus.Logger,
	redisRepo repository.CacheRepo, emailWriter *kafka.Writer,
	cfg *config.Config, metrics metrics.Metrics,
	profilesService profiles_service.ProfilesServiceV1Client) *AccountService {

	errorHandler := newErrorHandler(logger)
	return &AccountService{repo: repo,
		logger:                 logger,
		redisRepo:              redisRepo,
		nonActivatedAccountTTL: time.Hour,
		emailWriter:            emailWriter,
		cfg:                    cfg,
		metrics:                metrics,
		errorHandler:           errorHandler,
		profilesService:        profilesService,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context,
	in *accounts_service.CreateAccountRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.CreateAccount")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	vErr := validateSignupInput(in)
	if vErr != nil {
		err = vErr
		return nil, s.errorHandler.createExtendedErrorResponce(ErrFailedValidation, vErr.DeveloperMessage, vErr.UserMessage)
	}
	exist, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrInternal, err.Error(), "")
	}
	if exist {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrAlreadyExist, "", "a user with this email address already exists. "+
			"please try another one or simple log in")
	}

	inCache, err := s.redisRepo.RegistrationCache.IsAccountInCache(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if inCache {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrAlreadyExist, "", "a user with this email address already exists. "+
			"please try another one or verify email and log in")
	}

	s.logger.Info("Generating hash from password")
	password_hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), config.GetConfig().Crypto.BcryptCost)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, "can't generate hash")
	}

	err = s.redisRepo.RegistrationCache.CacheAccount(ctx, in.Email,
		repository.CachedAccount{Username: in.Username, Password: string(password_hash)}, s.nonActivatedAccountTTL)

	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error()+" can't cache account")
	}
	return &emptypb.Empty{}, nil
}

type emailData struct {
	Email    string  `json:"email"`
	URL      string  `json:"url"`
	MailType string  `json:"mail_type"`
	LinkTTL  float64 `json:"link_TTL"`
}

func (s *AccountService) RequestAccountVerificationToken(ctx context.Context,
	in *accounts_service.VerificationTokenRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"AccountService.RequestAccountVerificationToken")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	inAccountDB, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if inAccountDB {
		return nil, s.errorHandler.createErrorResponce(ErrAccountAlreadyActivated, "")
	}

	vErr := validateEmail(in.Email)
	if vErr != nil {
		err = vErr
		return nil, s.errorHandler.createExtendedErrorResponce(ErrInvalidArgument, vErr.DeveloperMessage, vErr.UserMessage)
	}

	inCache, err := s.redisRepo.RegistrationCache.IsAccountInCache(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if !inCache {
		s.metrics.IncCacheMiss("RequestAccountVerificationToken")
		return nil, s.errorHandler.createExtendedErrorResponce(ErrNotFound, "", "a account with this email address not exist")
	}
	s.metrics.IncCacheHits("RequestAccountVerificationToken")

	cfg := config.GetConfig()
	token, err := jwt.GenerateToken(in.Email, cfg.JWT.VerifyAccountToken.Secret, cfg.JWT.VerifyAccountToken.TTL)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	URL := fmt.Sprintf("%s/%s", in.URL, token)
	LinkTTL := cfg.JWT.VerifyAccountToken.TTL.Seconds()
	body, err := json.Marshal(emailData{Email: in.Email, URL: URL, MailType: "account/activation", LinkTTL: LinkTTL})
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
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
	in *accounts_service.VerifyAccountRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.VerifyAccount")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Parsing token")
	email, err := jwt.ParseToken(in.VerificationToken, config.GetConfig().JWT.VerifyAccountToken.Secret)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInvalidArgument, err.Error())
	}
	err = s.createAccountAndProfile(ctx, email)
	return &emptypb.Empty{}, err
}

func (s *AccountService) createAccountAndProfile(ctx context.Context, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.createAccountAndProfile")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	s.logger.Info("Checking account existing in cache")
	acc, err := s.redisRepo.RegistrationCache.GetCachedAccount(ctx, email)
	if errors.Is(err, redis.Nil) {
		s.metrics.IncCacheMiss("createAccountAndProfile")
		return s.errorHandler.createErrorResponce(ErrNotFound, err.Error())
	}
	if err != nil {
		s.metrics.IncCacheMiss("createAccountAndProfile")
		return s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	s.metrics.IncCacheHits("createAccountAndProfile")

	account := model.Account{
		Email:            email,
		Password:         string(acc.Password),
		RegistrationDate: time.Now().In(time.UTC).In(time.UTC),
	}

	s.logger.Info("Creating account")
	tx, accountID, err := s.repo.CreateAccount(ctx, account)
	if err != nil {
		return s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	defer tx.Rollback()
	profile := &profiles_service.CreateProfileRequest{
		AccountID:        accountID,
		Email:            email,
		Username:         acc.Username,
		RegistrationDate: timestamppb.New(account.RegistrationDate),
	}
	_, err = s.profilesService.CreateProfile(ctx, profile)
	if err != nil {
		s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	tx.Commit()

	//The error is not critical, the data will still be deleted from the cache.
	if err = s.redisRepo.RegistrationCache.DeleteAccountFromCache(ctx, email); err != nil {
		s.logger.Warning("can't delete account from registration cache: ", err.Error())
	}

	return nil
}

func (s *AccountService) SignIn(ctx context.Context,
	in *accounts_service.SignInRequest) (*accounts_service.AccessResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.SignIn")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	if net.ParseIP(in.ClientIP) == nil {
		err = s.errorHandler.createErrorResponce(ErrInvalidClientIP, "invalid client ip address")
		return nil, err
	}

	MachineID, err := s.getMachineIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Getting user by email")
	account, err := s.repo.GetAccountByEmail(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrNotFound, err.Error(), "account not found")
	}

	s.logger.Info("Password and hash comparison")
	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(in.Password)); err != nil {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrInvalidArgument, err.Error(), "invalid login or password")
	}

	s.logger.Info("Caching session")
	SessionID := uuid.NewString()
	if err = s.redisRepo.SessionsCache.CacheSession(ctx, model.SessionCache{SessionID: SessionID,
		AccountID: account.UUID, MachineID: MachineID, ClientIP: in.ClientIP, LastActivity: time.Now().In(time.UTC)}); err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return &accounts_service.AccessResponce{SessionID: SessionID}, nil
}

// --------------------- CONTEXTS ---------------------
const (
	AccountIdContext = "X-Account-Id"
	SessionIdContext = "X-Session-Id"
	MachineIdContext = "X-Machine-Id"
)

//-----------------------------------------------------

func (s *AccountService) GetAccountID(ctx context.Context,
	in *emptypb.Empty) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.GetAccountID")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	s.logger.Info("Checking session")
	cache, SessionID, _, err := s.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		s.logger.Info("Updating last activity for given session")
		span, ctx := opentracing.StartSpanFromContext(context.Background(),
			"AccountService.GetAccountID.UpdateLastActivityForSession")
		defer span.Finish()

		err := s.redisRepo.SessionsCache.UpdateLastActivityForSession(ctx, cache, SessionID, time.Now().In(time.UTC))
		if err != nil {
			s.logger.Warning("Session last activity not updated, error: ", err.Error())
		}

		span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	}()

	header := metadata.Pairs(AccountIdContext, cache.AccountID)
	grpc.SetHeader(ctx, header)
	return &emptypb.Empty{}, nil
}

func (s *AccountService) Logout(ctx context.Context,
	in *emptypb.Empty) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.Logout")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Checking session")
	cache, SessionID, _, err := s.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Terminating session: ", SessionID)
	err = s.redisRepo.SessionsCache.TerminateSessions(ctx, []string{SessionID}, cache.AccountID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) RequestChangePasswordToken(ctx context.Context,
	in *accounts_service.ChangePasswordTokenRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.RequestChangePasswordToken")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	exist, err := s.repo.IsAccountWithEmailExist(ctx, in.Email)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if !exist {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrNotFound, "", "a account with this email address not exist")
	}

	token, err := jwt.GenerateToken(in.Email, s.cfg.JWT.ChangePasswordToken.Secret, s.cfg.JWT.ChangePasswordToken.TTL)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	URL := fmt.Sprintf("%s/%s", in.URL, token)
	LinkTTL := s.cfg.JWT.ChangePasswordToken.TTL.Seconds()
	body, err := json.Marshal(emailData{Email: in.Email, URL: URL, MailType: "account/forget-password", LinkTTL: LinkTTL})
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
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
	in *accounts_service.ChangePasswordRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.ChangePassword")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Validating incoming password")
	vErr := validatePassword(in.ChangePasswordToken)
	if vErr != nil {
		err = vErr
		return nil, s.errorHandler.createExtendedErrorResponce(ErrInvalidArgument, vErr.DeveloperMessage, vErr.UserMessage)
	}

	s.logger.Info("Parsing jwt token")
	email, err := jwt.ParseToken(in.ChangePasswordToken, config.GetConfig().JWT.ChangePasswordToken.Secret)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInvalidArgument, err.Error())
	}

	s.logger.Info("Checking account existing in DB")
	exist, err := s.repo.IsAccountWithEmailExist(ctx, email)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if !exist {
		return nil, s.errorHandler.createExtendedErrorResponce(ErrNotFound, "", "account not found")
	}

	GeneratingHashSpan, _ := opentracing.StartSpanFromContext(ctx,
		"AccountService.ChangePassword.GenerateHash")

	s.logger.Info("Generating hash for incoming password")
	password_hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), config.GetConfig().Crypto.BcryptCost)
	if err != nil {
		GeneratingHashSpan.Finish()
		return nil, s.errorHandler.createErrorResponce(ErrInternal, "can't generate hash.")
	}
	GeneratingHashSpan.Finish()

	s.logger.Info("Changing account password")
	err = s.repo.ChangePassword(ctx, email, string(password_hash))
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) GetAllSessions(ctx context.Context,
	in *emptypb.Empty) (*accounts_service.AllSessionsResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.GetAllSessions")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Checking session")
	cache, _, _, err := s.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Getting sessions for account ", cache.AccountID)
	sessions, err := s.redisRepo.SessionsCache.GetSessionsForAccount(ctx, cache.AccountID)
	if err != nil && err != redis.Nil {
		s.metrics.IncCacheMiss("GetAllSessions")
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	s.metrics.IncCacheHits("GetAllSessions")

	s.logger.Info("Converting cache data into responce")
	sessionsInfo := make(map[string]*accounts_service.SessionInfo, len(sessions))
	for key, session := range sessions {
		sessionsInfo[key] = &accounts_service.SessionInfo{
			ClientIP:     session.ClientIP,
			MachineID:    session.MachineID,
			LastActivity: timestamppb.New(session.LastActivity),
		}
	}
	s.logger.Info("Cache data successfully converted into responce, sending responce")

	return &accounts_service.AllSessionsResponce{Sessions: sessionsInfo}, nil
}

func (s *AccountService) TerminateSessions(ctx context.Context,
	in *accounts_service.TerminateSessionsRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.TerminateSessions")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Checking session")
	cache, _, _, err := s.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Terminating sessions")
	if err = s.redisRepo.SessionsCache.TerminateSessions(ctx, in.SessionsToTerminate, cache.AccountID); err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context,
	in *emptypb.Empty) (*emptypb.Empty, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.DeleteAccount")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Checking session")
	cache, _, _, err := s.checkSession(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := s.repo.DeleteAccount(ctx, cache.AccountID)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	defer tx.Rollback()
	_, err = s.profilesService.DeleteProfile(ctx, &profiles_service.DeleteProfileRequest{AccountID: cache.AccountID})
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return &emptypb.Empty{}, nil
}

func (s *AccountService) checkSession(ctx context.Context) (cache model.SessionCache, SessionID, ClientIP string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AccountService.checkSession")
	defer span.Finish()
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Getting session id from ctx")
	SessionID, err = s.getSessionIDFromCtx(ctx)
	if err != nil {
		return
	}

	s.logger.Info("Getting client ip from ctx")
	MachineID, err := s.getMachineIDFromCtx(ctx)
	if err != nil {
		return
	}

	s.logger.Info("Getting session cache")
	cache, err = s.redisRepo.SessionsCache.GetSessionCache(ctx, SessionID)
	if errors.Is(err, repository.ErrSessionNotFound) {
		s.metrics.IncCacheMiss("checkSession")
		err = s.errorHandler.createErrorResponce(ErrSessisonNotFound, "")
		return
	} else if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error())
		return
	}

	s.metrics.IncCacheHits("checkSession")

	if MachineID != cache.MachineID {
		err = s.errorHandler.createErrorResponce(ErrAccessDenied, "")
		return
	}

	return
}

func (s *AccountService) getSessionIDFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", s.errorHandler.createErrorResponce(ErrNoCtxMetaData, "")
	}

	sessionID := md.Get(SessionIdContext)
	if len(sessionID) == 0 || sessionID[0] == "" {
		return "", s.errorHandler.createErrorResponce(ErrInvalidSessionId, "no session id provided")
	}

	return sessionID[0], nil
}

// func (s *AccountService) getClientIPFromCtx(ctx context.Context) (string, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return "", s.errorHandler.createErrorResponce(ErrNoCtxMetaData, "no context metadata provided")
// 	}

// 	clientIP := md.Get(ClientIpContext)
// 	if len(clientIP) == 0 || clientIP[0] == "" {
// 		return "", s.errorHandler.createErrorResponce(ErrInvalidClientIP, "no client ip provided")
// 	}

// 	if net.ParseIP(clientIP[0]) == nil {
// 		return "", s.errorHandler.createErrorResponce(ErrInvalidClientIP, "invalid client ip address")
// 	}

// 	return clientIP[0], nil
// }

func (s *AccountService) getMachineIDFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", s.errorHandler.createErrorResponce(ErrNoCtxMetaData, "no context metadata provided")
	}

	MachineID := md.Get(MachineIdContext)
	if len(MachineID) == 0 || MachineID[0] == "" {
		return "", s.errorHandler.createErrorResponce(ErrInvalidMachineID, "no client ip provided")
	}

	return MachineID[0], nil
}
