package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/service"
	account_service "github.com/Falokut/online_cinema_ticket_office/account_service/pkg/account_service/protos"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
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

func (s *server) Run(cfg Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	account_service.RegisterAccountServiceV1Server(s.grpcServer, s.service)
	return s.grpcServer.Serve(lis)
}

func (s *server) ShutDown() {
	s.logger.Println("Shutting down")
	s.service.ShutDown()
	s.grpcServer.GracefulStop()
	wg.Wait()
}
