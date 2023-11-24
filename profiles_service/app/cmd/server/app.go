package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/service"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_resizer"
	image_storage_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_storage_service/v1/protos"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/metrics"
	profiles_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/profiles_service/v1/protos"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sirupsen/logrus"
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

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}

	logger.Info("Repository initializing")
	repo := repository.NewProfileRepository(database)
	defer repo.Shutdown()

	logger.Info("GRPC Client initializing")
	conn, err := getImageStorageConnection(appCfg)
	if err != nil {
		logger.Fatal(err)
	}
	imageStorageService := image_storage_service.NewImagesStorageServiceV1Client(conn)
	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{database}, appCfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Fatalf("Shutting down, can't run healthcheck endpoint %s", err.Error())
		}
	}()
	logger.Info("Service initializing")
	images_resizer.SetLogger(logging.GetLogger().Logger)
	imagesService := service.NewImageService(getImageStorageConfig(appCfg),
		logger, imageStorageService)

	service := service.NewProfilesService(repo, logger.Logger, metric, imagesService)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	s.Run(getListenServerConfig(appCfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
}
func getImageStorageConnection(cfg *config.Config) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg.ImageService.StorageAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
}

func getImageStorageConfig(cfg *config.Config) service.ImageServiceConfig {
	return service.ImageServiceConfig{
		ImageWidth:             cfg.ImageService.ImageWidth,
		ImageHeight:            cfg.ImageService.ImageHeight,
		ImageResizeType:        cfg.ImageService.ImageResizeType,
		ImageResizeMethod:      images_resizer.ResolveResizeMethod(cfg.ImageService.ImageResizeMethod),
		BaseProfilePictureUrl:  cfg.ImageService.BaseProfilePictureUrl,
		ProfilePictureCategory: cfg.ImageService.ProfilePictureCategory,
		MaxImageWidth:          cfg.ImageService.MaxImageWidth,
		MaxImageHeight:         cfg.ImageService.MaxImageHeight,
		MinImageWidth:          cfg.ImageService.MinImageWidth,
		MinImageHeight:         cfg.ImageService.MinImageHeight,
	}
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:           cfg.Listen.Mode,
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
		ServiceDesc:    &profiles_service.ProfilesServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(profiles_service.ProfilesServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return profiles_service.RegisterProfilesServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}
