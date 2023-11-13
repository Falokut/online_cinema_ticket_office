package interceptors

import (
	"context"
	"net/http"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/metrics"
	"github.com/felixge/httpsnoop"
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
		status = grpc_errors.ConvertGrpcCodeIntoHTTP(grpc_errors.GetGrpcCode(err))
	}
	im.metr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.metr.IncHits(status, info.FullMethod, info.FullMethod)

	return resp, err
}

func (im *InterceptorManager) RestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		im.logger.Infof("Method: %s, Path: %s, Time: %v, status code: %v", request.Method, request.URL.Path,
			m.Duration, m.Code)
	})
}

func (im *InterceptorManager) RestMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)

		status := m.Code
		im.metr.ObserveResponseTime(status, request.Method, request.URL.Path, m.Duration.Seconds())
		im.metr.IncHits(status, request.Method, request.URL.Path)
	})
}
