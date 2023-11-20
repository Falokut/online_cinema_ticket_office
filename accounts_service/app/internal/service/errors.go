package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound                = errors.New("not found")
	ErrNoCtxMetaData           = errors.New("no context metadata")
	ErrInvalidSessionId        = errors.New("invalid session id")
	ErrAlreadyExist            = errors.New("already exist")
	ErrInvalidClientIP         = errors.New("invalid client ip")
	ErrAccessDenied            = errors.New("access denied. Invalid session or client ip")
	ErrInternal                = errors.New("internal error")
	ErrAccountAlreadyActivated = errors.New("account already activated")
	ErrInvalidArgument         = errors.New("invalid input data")
	ErrFailedValidation        = errors.New("validation failed")
	ErrSessisonNotFound        = errors.New("session with specified id not found")
)

var errorCodes = map[error]codes.Code{
	redis.Nil:                  codes.NotFound,
	ErrNotFound:                codes.NotFound,
	ErrNoCtxMetaData:           codes.Unauthenticated,
	ErrInvalidSessionId:        codes.Unauthenticated,
	ErrSessisonNotFound:        codes.Unauthenticated,
	ErrAlreadyExist:            codes.AlreadyExists,
	ErrInvalidClientIP:         codes.InvalidArgument,
	ErrFailedValidation:        codes.InvalidArgument,
	ErrAccessDenied:            codes.PermissionDenied,
	ErrInternal:                codes.Internal,
	ErrAccountAlreadyActivated: codes.AlreadyExists,
}

type errorHandler struct {
	logger logging.Logger
}

func newErrorHandler(logger logging.Logger) errorHandler {
	return errorHandler{
		logger: logger,
	}
}

func (e *errorHandler) createErrorResponce(err error, developerMessage string) error {
	var msg string
	if len(developerMessage) == 0 {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", developerMessage, err)
	}

	err = status.Error(grpc_errors.GetGrpcCode(err), msg)
	e.logger.Error(err)
	return err
}

func (e *errorHandler) createExtendedErrorResponce(err error, DeveloperMessage, UserMessage string) error {
	var msg string
	if DeveloperMessage == "" {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", DeveloperMessage, err)
	}

	extErr := status.New(grpc_errors.GetGrpcCode(err), msg)
	if len(UserMessage) > 0 {
		extErr, _ = extErr.WithDetails(&accounts_service.UserErrorMessage{Message: UserMessage})
		if extErr == nil {
			e.logger.Error(err)
			return err
		}
	}

	e.logger.Error(extErr)
	return extErr.Err()
}

func init() {
	grpc_errors.RegisterErrors(errorCodes)
}
