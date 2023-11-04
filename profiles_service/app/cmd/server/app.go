package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/service"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/metrics"
	"github.com/opentracing/opentracing-go"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	log_level, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(log_level)

	logger.Println(appCfg.JaegerConfig)
	tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}

	logger.Info("Repository initializing")
	repo := repository.NewProfileRepository(database)

	logger.Info("Service initializing")
	service := service.NewProfileService(repo, logger)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(appCfg.PrometheusConfig.Address,
		appCfg.PrometheusConfig.Name, logger)
	if err != nil {
		logger.Fatal(err)
	}

	s := server.NewServer(logger, service)
	logger.Info("Server initializing")
	go func() {
		logger.Info("Server running")
		if err := s.Run(getServerConfig(appCfg), metric); err != nil {
			logger.Fatalf("%s", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	if err := repo.Shutdown(); err != nil {
		logger.Error("Error when closing connection with db: ", err)
	}
	s.ShutDown()
}

func getServerConfig(appCfg *config.Config) server.Config {
	return server.Config{
		Address: appCfg.Listen.BindIP,
		Port:    appCfg.Listen.Port}
}
