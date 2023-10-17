package server

import (
	"fmt"
	"net"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/service"
	profile_service "github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/profile_service/protos"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

type Config struct {
	Address string
	Port    uint16
}

func (s *Server) Run(cfg Config, service *service.ProfileService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		return err
	}
	servReg := grpc.NewServer()
	profile_service.RegisterProfileServiceV1Server(servReg, service)
	servReg.GracefulStop()

	return servReg.Serve(lis)
}

func (s *Server) Shutdown() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
