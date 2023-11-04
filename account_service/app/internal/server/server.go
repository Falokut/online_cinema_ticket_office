package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/interceptors"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/service"
	account_service "github.com/Falokut/online_cinema_ticket_office/account_service/pkg/account_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/metrics"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

type server struct {
	service    *service.AccountService
	logger     logging.Logger
	grpcServer *grpc.Server
}

func NewServer(logger logging.Logger, service *service.AccountService) server {
	return server{logger: logger, service: service}
}

type Config struct {
	Address string
	Port    uint16
}

func (s *server) Run(cfg Config, metric metrics.Metrics) error {
	im := interceptors.NewInterceptorManager(s.logger, metric)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(im.Logger), grpc.ChainUnaryInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	))

	grpc_prometheus.Register(s.grpcServer)
	http.Handle("/metrics", promhttp.Handler())

	account_service.RegisterAccountServiceV1Server(s.grpcServer, s.service)
	return s.grpcServer.Serve(lis)
}

func (s *server) ShutDown() {
	s.logger.Println("Shutting down")
	s.grpcServer.GracefulStop()
	wg.Wait()
}
