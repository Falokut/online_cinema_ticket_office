package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
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
)

var errorCodes = map[error]codes.Code{
	ErrProfileNotFound:          codes.NotFound,
	ErrNoCtxMetaData:            codes.Unauthenticated,
	ErrInvalidAccountId:         codes.Unauthenticated,
	ErrCantUpdateProfilePicture: codes.Internal,
	ErrInternal:                 codes.Internal,
}

type errorHandler struct {
	logger *logrus.Logger
}

func newErrorHandler(logger *logrus.Logger) errorHandler {
	return errorHandler{
		logger: logger,
	}
}

func (e *errorHandler) createErrorResponce(err error, errorMessage string) error {
	var msg string
	if errorMessage == "" {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", errorMessage, err)
	}

	e.logger.Error(status.Error(grpc_errors.GetGrpcCode(err), msg))
	return status.Error(grpc_errors.GetGrpcCode(err), err.Error())
}

func init() {
	grpc_errors.RegisterErrors(errorCodes)
}
