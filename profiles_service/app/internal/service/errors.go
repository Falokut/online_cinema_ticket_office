package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
	profiles_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/profiles_service/v1/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrProfileNotFound          = errors.New("profile not found")
	ErrCantUpdateProfilePicture = errors.New("can't update profile picture")
	ErrInternal                 = errors.New("internal")
	ErrNoCtxMetaData            = errors.New("no context metadata")
	ErrInvalidAccountId         = errors.New("invalid account id")
	ErrInvalidImage             = errors.New("invalid image")
)

var errorCodes = map[error]codes.Code{
	ErrProfileNotFound:          codes.NotFound,
	ErrNoCtxMetaData:            codes.Unauthenticated,
	ErrInvalidAccountId:         codes.Unauthenticated,
	ErrCantUpdateProfilePicture: codes.Internal,
	ErrInternal:                 codes.Internal,
	ErrInvalidImage:             codes.InvalidArgument,
}

type errorHandler struct {
	logger *logrus.Logger
}

func newErrorHandler(logger *logrus.Logger) errorHandler {
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
		extErr, _ = extErr.WithDetails(&profiles_service.UserErrorMessage{Message: UserMessage})
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
