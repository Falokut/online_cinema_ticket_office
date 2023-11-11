package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCantGetImageByID    = errors.New("can't get image with specified id")
	ErrCantFindImageByID   = errors.New("can't find image with specified id")
	ErrUnsupportedFileType = errors.New("the received file type is not supported")
	ErrZeroSizeFile        = errors.New("the received file has zero size")
	ErrImageTooLarge       = errors.New("image is too large")
	ErrCantWriteChunkData  = errors.New("can't write chunk data")
	ErrCantReplaceImage    = errors.New("can't replace image")
	ErrReceivedNilRequest  = errors.New("the received request is nil")
)

var errorCodes = map[error]codes.Code{
	ErrCantGetImageByID:    codes.Internal,
	ErrCantFindImageByID:   codes.NotFound,
	ErrUnsupportedFileType: codes.InvalidArgument,
	ErrZeroSizeFile:        codes.InvalidArgument,
	ErrImageTooLarge:       codes.InvalidArgument,
	ErrCantWriteChunkData:  codes.Internal,
	ErrCantReplaceImage:    codes.Internal,
	ErrReceivedNilRequest:  codes.InvalidArgument,
}

type errorHandler struct {
	logger logging.Logger
}

func newErrorHandler(logger logging.Logger) errorHandler {
	grpc_errors.RegisterErrors(errorCodes)
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

	responceErr := status.Error(grpc_errors.GetGrpcCode(err), msg)
	e.logger.Error(responceErr)
	return responceErr
}
