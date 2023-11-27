package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/service"
	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/metrics"
	profiles_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/profiles_service/v1/protos"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

func main() {
	logging.NewEntry(logging.FileAndConsoleOutput)
	logger := logging.GetLogger()
	// Test
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
		logger.Infof("Metrics server running at %s:%s", appCfg.PrometheusConfig.ServerConfig.Host,
			appCfg.PrometheusConfig.ServerConfig.Port)
		if err := metrics.RunMetricServer(appCfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info("Registration cache initializing")
	regCacheOpt := appCfg.RegistrationCacheOptions.ConvertToRedisOptions()
	registrationCache, err := repository.NewRedisRegistrationCache(regCacheOpt, logger)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis registration cache is not established: %s options: %v",
			err.Error(), regCacheOpt)
	}
	defer registrationCache.Shutdown()

	logger.Info("Sessions cache initializing")
	sessionCacheOpt := appCfg.SessionCacheOptions.ConvertToRedisOptions()
	accountSessionCacheOpt := appCfg.AccountSessionsCacheOptions.ConvertToRedisOptions()
	sessionsCache, err := repository.NewSessionCache(sessionCacheOpt,
		accountSessionCacheOpt, logger, appCfg.SessionsTTL)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis token cache is not established: %s options: %v",
			err.Error(), sessionCacheOpt)
	}
	defer sessionsCache.Shutdown()

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}
	defer database.Close()

	redisRepo := repository.NewCacheRepository(registrationCache, sessionsCache)

	logger.Info("Repository initializing")
	repo := repository.NewAccountRepository(database)

	kafkaWriter := getKafkaWriter(appCfg.EmailKafka)
	defer kafkaWriter.Close()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{database, registrationCache, sessionsCache}, appCfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Fatalf("Shutting down, can't run healthcheck endpoint %s", err.Error())
		}
	}()
	logger.Info("Healthcheck initialized")

	cc, err := getProfilesServiceConnection(appCfg)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the profiles service is not established: %s", err.Error())
	}
	profilesService := profiles_service.NewProfilesServiceV1Client(cc)
	logger.Info("Service initializing")
	service := service.NewAccountService(repo,
		logger, redisRepo, kafkaWriter, appCfg, metric, profilesService)

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
		Mode:                   cfg.Listen.Mode,
		Host:                   cfg.Listen.Host,
		Port:                   cfg.Listen.Port,
		AllowedHeaders:         cfg.Listen.AllowedHeaders,
		AllowedOutgoingHeaders: cfg.Listen.AllowedOutgoingHeaders,
		ServiceDesc:            &accounts_service.AccountsServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(accounts_service.AccountsServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return accounts_service.RegisterAccountsServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}

func getProfilesServiceConnection(cfg *config.Config) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg.ProfilesServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
}

func getKafkaWriter(cfg config.KafkaConfig) *kafka.Writer {
	w := &kafka.Writer{
		Addr:   kafka.TCP(cfg.Brokers...),
		Topic:  cfg.Topic,
		Logger: logging.GetLogger().Logger,
	}
	return w
}
