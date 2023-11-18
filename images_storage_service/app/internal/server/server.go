package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Falokut/interceptors"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/service"
	images_storage_service "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/metrics"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	service                *service.ImagesStorageService
	logger                 *logrus.Logger
	grpcServer             *grpc.Server
	AllowedHeaders         []string
	AllowedOutgoingHeaders map[string]string
	im                     *interceptors.InterceptorManager
	mux                    cmux.CMux
}

type Config struct {
	Host                   string
	Port                   string
	Mode                   string
	AllowedHeaders         []string
	AllowedOutgoingHeaders map[string]string
}

func NewServer(logger *logrus.Logger, service *service.ImagesStorageService) server {
	return server{logger: logger, service: service}
}

func (s *server) Run(cfg Config, metric metrics.Metrics) {
	Mode := strings.ToUpper(cfg.Mode)
	s.logger.Info("start running server on mode: " + Mode)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		s.logger.Fatal("error while listening", err)
	}
	s.mux = cmux.New(lis)
	s.im = interceptors.NewInterceptorManager(s.logger, metric)

	switch Mode {
	case "REST":
		s.RunRestAPI(cfg, metric)
	case "GRPC":
		s.RunGRPC(cfg, metric)
	case "BOTH":
		s.RunRestAPI(cfg, metric)
		s.RunGRPC(cfg, metric)
	}
	if err := s.mux.Serve(); err != nil {
		s.logger.Fatal(err)
	}
	grpc_prometheus.Register(s.grpcServer)

	s.logger.Info("server running on mode: " + Mode)
}

func (s *server) RunGRPC(cfg Config, metric metrics.Metrics) {
	s.logger.Info("GRPC server initializing")

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(s.im.Logger), grpc.ChainUnaryInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(),
		s.im.Metrics,
		grpcrecovery.UnaryServerInterceptor(),
	), grpc.StreamInterceptor(s.im.StreamLogger), grpc.ChainStreamInterceptor(
		grpc_ctxtags.StreamServerInterceptor(),
		s.im.StreamMetrics,
		grpcrecovery.StreamServerInterceptor(),
	), grpc.Creds(insecure.NewCredentials()),
	)

	images_storage_service.RegisterImagesStorageServiceV1Server(s.grpcServer, s.service)
	go func() {
		grpcL := s.mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		if err := s.grpcServer.Serve(grpcL); err != nil {
			s.logger.Fatalf("GRPC error while serving: %v", err)
		}
	}()
	s.logger.Infof("GRPC server initialized. Listen on %s:%s", cfg.Host, cfg.Port)

}

func (s *server) RunRestAPI(cfg Config, metric metrics.Metrics) {
	s.logger.Info("REST server initializing")

	s.AllowedHeaders = cfg.AllowedHeaders
	s.AllowedOutgoingHeaders = cfg.AllowedOutgoingHeaders

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(s.headerMatcherFunc),
	)

	if err := images_storage_service.RegisterImagesStorageServiceV1HandlerServer(context.Background(),
		mux, s.service); err != nil {
		s.logger.Fatalf("REST server error while registering handler server: %v", err)
	}

	server := http.Server{
		Handler: s.im.RestLogger(s.im.RestMetrics(mux)),
	}

	s.logger.Info("Rest server initializing")
	go func() {
		restL := s.mux.Match(cmux.HTTP1Fast())

		if err := server.Serve(restL); err != nil {
			s.logger.Fatalf("REST server error while serving: %v", err)
		}
	}()

	s.logger.Infof("REST server initialized. Listen on %s:%s", cfg.Host, cfg.Port)
}

func (s *server) Shutdown() {
	s.logger.Info("shutting down")
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	s.mux.Close()
}

func (s *server) headerMatcherFunc(header string) (string, bool) {
	s.logger.Debugf("Received %s header", header)
	for _, AllowedHeader := range s.AllowedHeaders {
		if header == AllowedHeader {
			return header, true
		}
	}

	return runtime.DefaultHeaderMatcher(header)
}
