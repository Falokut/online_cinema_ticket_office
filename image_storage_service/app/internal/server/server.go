package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/handlers"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/interceptors"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/service"
	image_storage_service "github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/image_storage_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/metrics"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type server struct {
	service        *service.ImageStorageService
	logger         logging.Logger
	grpcServer     *grpc.Server
	AllowedHeaders []string
}

type Config struct {
	Host           string
	Port           string
	AllowedHeaders []string
}

func NewServer(logger logging.Logger, service *service.ImageStorageService) server {
	return server{logger: logger, service: service}
}

func (s *server) Run(cfg Config, metric metrics.Metrics) error {
	s.logger.Info("Start running server")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}
	m := cmux.New(lis)

	if err := s.RunRestGRPCServer(m, cfg, metric); err != nil {
		return err
	}

	s.RunRestAPIServer(m, cfg)

	if err := m.Serve(); err != nil {
		return err
	}

	s.logger.Info("Server running")

	return nil
}
func (s *server) RunRestGRPCServer(m cmux.CMux, cfg Config, metric metrics.Metrics) error {
	s.logger.Info("GRPC server initializing")
	im := interceptors.NewInterceptorManager(s.logger, metric)

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(im.Logger), grpc.ChainUnaryInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	))

	grpc_prometheus.Register(s.grpcServer)

	image_storage_service.RegisterImageStorageServiceV1Server(s.grpcServer, s.service)
	s.logger.Infof("GRPC server initialized. Start listening on %s:%s", cfg.Host, cfg.Port)
	go func() {
		grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		if err := s.grpcServer.Serve(grpcL); err != nil {
			s.logger.Fatalf("GRPC error while serving: %v", err)
		}
	}()

	return nil
}
func (s *server) RunRestAPIServer(m cmux.CMux, cfg Config) {
	s.logger.Info("Rest server initializing")
	h := handlers.NewHandler(s.logger, s.service)
	rest_m := h.RegisterHandler()

	go func() {
		restL := m.Match(cmux.HTTP1Fast())
		s.logger.Infof("REST server initialized. Listen on %s:%s", cfg.Host, cfg.Port)
		if err := http.Serve(restL, rest_m); err != nil {
			s.logger.Fatalf("REST server error while serving: %v", err)
		}
	}()
}
func (s *server) ShutDown() {
	s.logger.Println("Shutting down")
	s.grpcServer.GracefulStop()
}

func (s *server) headerMatcherFunc(header string) (string, bool) {
	for _, AllowedHeader := range s.AllowedHeaders {
		if strings.ToLower(header) == AllowedHeader {
			return AllowedHeader, true
		}
	}

	return runtime.DefaultHeaderMatcher(header)
}
