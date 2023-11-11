package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/internal/service"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/metrics"
	"github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
)

func main() {
	logging.NewEntry(logging.FileAndConsoleOutput)
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	log_level, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(log_level)
	if appCfg.MaxImageSize <= 0 {
		logger.Fatal("Max image size less or equal zero")
	}
	tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(appCfg.PrometheusConfig.Name)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(appCfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info("Local storage initializing")
	storage := repository.NewLocalStorage(logger, appCfg.BaseLocalStoragePath)

	logger.Info("Service initializing")
	service := service.NewImageStorageService(logger,
		service.Config{MaxImageSize: appCfg.MaxImageSize}, storage)

	logger.Info("Server initializing")
	s := server.NewServer(logger, service)
	s.Run(getListenServerConfig(appCfg), metric)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.ShutDown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
	}
}
