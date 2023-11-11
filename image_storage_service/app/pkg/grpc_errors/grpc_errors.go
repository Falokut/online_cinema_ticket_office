package grpc_errors

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
)

var errorCodes = map[error]codes.Code{
	sql.ErrNoRows:            codes.NotFound,
	context.Canceled:         codes.Canceled,
	context.DeadlineExceeded: codes.DeadlineExceeded,
	os.ErrNotExist:           codes.NotFound,
	os.ErrInvalid:            codes.Internal,
	os.ErrPermission:         codes.Internal,
	os.ErrDeadlineExceeded:   codes.DeadlineExceeded,
}

func RegisterErrors(addErrors map[error]codes.Code) {
	for key, val := range addErrors {
		errorCodes[key] = val
	}
}

// Map GRPC errors codes to http status
func ConvertGrpcCodeIntoHTTP(code codes.Code) int {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.AlreadyExists:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func GetGrpcCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	switch {
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	}

	code, ok := errorCodes[err]
	if !ok {
		return codes.Unknown
	}
	return code
}
