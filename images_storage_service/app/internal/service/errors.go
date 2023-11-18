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
	ErrCantFindImageByID   = errors.New("can't find image with specified id")
	ErrUnsupportedFileType = errors.New("the received file type is not supported")
	ErrZeroSizeFile        = errors.New("the received file has zero size")
	ErrImageTooLarge       = errors.New("image is too large")
	ErrCantWriteChunkData  = errors.New("can't write chunk data")
	ErrCantReplaceImage    = errors.New("can't replace image")
	ErrReceivedNilRequest  = errors.New("the received request is nil")
	ErrCantSaveImage       = errors.New("сan't save image to the storage")
	ErrCantDeleteImage     = errors.New("can't delete image")
)

var errorCodes = map[error]codes.Code{
	ErrCantFindImageByID:   codes.NotFound,
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
