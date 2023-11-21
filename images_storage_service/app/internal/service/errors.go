package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
	img_storage_serv "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCantFindImageByID   = errors.New("can't find image with specified id")
	ErrUnsupportedFileType = errors.New("the received file type is not supported")
	ErrZeroSizeFile        = errors.New("the received file has zero size")
	ErrImageTooLarge       = errors.New("image is too large")
	ErrCantWriteChunkData  = errors.New("can't write chunk data")
	ErrCantReplaceImage    = errors.New("can't replace image")
	ErrReceivedNilRequest  = errors.New("the received request is nil")
	ErrCantSaveImage       = errors.New("сan't save image to the storage")
	ErrCantDeleteImage     = errors.New("can't delete image")
	ErrInternal            = errors.New("internal error")
)

var errorCodes = map[error]codes.Code{
	ErrCantFindImageByID:   codes.NotFound,
	ErrInternal:            codes.Internal,
	ErrCantDeleteImage:     codes.Internal,
	ErrCantSaveImage:       codes.Internal,
	ErrUnsupportedFileType: codes.InvalidArgument,
	ErrZeroSizeFile:        codes.InvalidArgument,
	ErrImageTooLarge:       codes.InvalidArgument,
	ErrCantWriteChunkData:  codes.Internal,
	ErrCantReplaceImage:    codes.Internal,
	ErrReceivedNilRequest:  codes.InvalidArgument,
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
		extErr, _ = extErr.WithDetails(&img_storage_serv.UserErrorMessage{Message: UserMessage})
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
