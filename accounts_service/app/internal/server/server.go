package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/interceptors"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/service"
	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/metrics"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

type server struct {
	service        *service.AccountService
	logger         logging.Logger
	grpcServer     *grpc.Server
	AllowedHeaders []string
	mux            cmux.CMux
}

type Config struct {
	Host           string
	Port           string
	Mode           string
	AllowedHeaders []string
}

func NewServer(logger logging.Logger, service *service.AccountService) server {
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

	switch Mode {
	case "REST":
		s.RunRestAPI(cfg)
	case "GRPC":
		s.RunGRPC(cfg, metric)
	case "BOTH":
		s.RunRestAPI(cfg)
		s.RunGRPC(cfg, metric)
	}
	if err := s.mux.Serve(); err != nil {
		s.logger.Fatal(err)
	}

	s.logger.Info("server running on mode: " + Mode)
}

func (s *server) RunGRPC(cfg Config, metric metrics.Metrics) {
	s.logger.Info("GRPC server initializing")
	im := interceptors.NewInterceptorManager(s.logger, metric)

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(im.Logger), grpc.ChainUnaryInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	))

	grpc_prometheus.Register(s.grpcServer)
	accounts_service.RegisterAccountsServiceV1Server(s.grpcServer, s.service)
	go func() {
		grpcL := s.mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		if err := s.grpcServer.Serve(grpcL); err != nil {
			s.logger.Fatalf("GRPC error while serving: %v", err)
		}
	}()
	s.logger.Infof("GRPC server initialized. Listen on %s:%s", cfg.Host, cfg.Port)

}

func (s *server) RunRestAPI(cfg Config) {
	s.logger.Info("REST server initializing")

	s.AllowedHeaders = cfg.AllowedHeaders

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(s.headerMatcherFunc),
	)

	if err := accounts_service.RegisterAccountsServiceV1HandlerServer(context.Background(),
		mux, s.service); err != nil {
		s.logger.Fatalf("REST server error while registering handler server: %v", err)
	}

	s.logger.Info("Rest server initializing")
	go func() {
		restL := s.mux.Match(cmux.HTTP1Fast())
		if err := http.Serve(restL, mux); err != nil {
			s.logger.Fatalf("REST server error while serving: %v", err)
		}
	}()

	s.logger.Infof("REST server initialized. Listen on %s:%s", cfg.Host, cfg.Port)
}

func (s *server) ShutDown() {
	s.logger.Info("shutting down")
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	s.mux.Close()
	wg.Wait()
}

func (s *server) headerMatcherFunc(header string) (string, bool) {
	for _, AllowedHeader := range s.AllowedHeaders {
		if strings.ToLower(header) == AllowedHeader {
			return AllowedHeader, true
		}
	}

	return runtime.DefaultHeaderMatcher(header)
}
