package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/shutdown"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handlers http.Handler, logger logging.Logger, cfg *config.Config) error {
	s.httpServer = &http.Server{
		Addr:           ":" + cfg.Listen.Port,
		Handler:        handlers,
		MaxHeaderBytes: 1 << 20, //1MB
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
	}

	s.httpServer.RegisterOnShutdown(s.Shutdown)
	return s.httpServer.ListenAndServe()
}

func RunServer(handlers http.Handler, logger logging.Logger, cfg *config.Config) error {
	server := &http.Server{
		Handler:      handlers,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	var listener net.Listener
	var err error
	if cfg.Listen.Type == "sock" {
		listener, err = configureListenTypeSocket(logger, cfg)
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	}

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	s.httpServer.Close()
	shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		s.httpServer)
}

func configureListenTypeSocket(logger logging.Logger, cfg *config.Config) (net.Listener, error) {
	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	socketPath := path.Join(appDir, "app.sock")
	logger.Infof("socket path: %s", socketPath)

	logger.Info("create and listen unix socket")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
