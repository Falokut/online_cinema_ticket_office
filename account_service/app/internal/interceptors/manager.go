package interceptors

import (
	"context"
	"net/http"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InterceptorManager
type InterceptorManager struct {
	logger logging.Logger
	metr   metrics.Metrics
}

// InterceptorManager constructor
func NewInterceptorManager(logger logging.Logger, metr metrics.Metrics) *InterceptorManager {
	return &InterceptorManager{logger: logger, metr: metr}
}

// Logger Interceptor
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}

func (im *InterceptorManager) Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	var status = http.StatusOK
	if err != nil {
		status = grpc_errors.MapGRPCErrCodeToHttpStatus(grpc_errors.ParseGRPCErrStatusCode(err))
	}
	im.metr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.metr.IncHits(status, info.FullMethod, info.FullMethod)

	return resp, err
}
