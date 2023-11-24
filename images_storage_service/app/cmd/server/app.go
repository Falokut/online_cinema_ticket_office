package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/service"
	img_storage_serv "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/images_storage_service/v1/protos"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/metrics"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	storage := repository.NewLocalStorage(logger.Logger, appCfg.BaseLocalStoragePath)
	defer storage.Shutdown()
	logger.Info("Service initializing")
	service := service.NewImagesStorageService(logger.Logger,
		service.Config{MaxImageSize: appCfg.MaxImageSize}, storage, metric)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	s.Run(getListenServerConfig(appCfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:           cfg.Listen.Mode,
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
		ServiceDesc:    &img_storage_serv.ImagesStorageServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(img_storage_serv.ImagesStorageServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return img_storage_serv.RegisterImagesStorageServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}
