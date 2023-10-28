package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/service"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/account_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/metrics"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	log_level, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println(appCfg.JaegerConfig)
	tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Logger.SetLevel(log_level)
	logger.Info("Registration cache initializing")
	regCacheOpt := appCfg.RegistrationCacheOptions.ConvertToRedisOptions()
	registrationCache, err := repository.NewRedisRegistrationCache(regCacheOpt, logger)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis registration cache is not established: %s options: %v",
			err.Error(), regCacheOpt)
	}

	logger.Info("Token cache initializing")
	sessionCacheOpt := appCfg.SessionCacheOptions.ConvertToRedisOptions()
	accountSessionCacheOpt := appCfg.AccountSessionsCacheOptions.ConvertToRedisOptions()
	sessionCache, err := repository.NewSessionCache(sessionCacheOpt,
		accountSessionCacheOpt, logger, appCfg.SessionsTTL)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis token cache is not established: %s options: %v",
			err.Error(), sessionCacheOpt)
	}

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}

	redisRepo := repository.NewCacheRepository(registrationCache, sessionCache)

	logger.Info("Repository initializing")
	repo := repository.NewAccountRepository(database)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(appCfg.PrometheusConfig.Address,
		appCfg.PrometheusConfig.Name, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Service initializing")
	service := service.NewAccountService(repo,
		logger, redisRepo, GetKafkaWriter(appCfg.EmailKafka), appCfg, metric)

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
	s.ShutDown()
}

func getServerConfig(appCfg *config.Config) server.Config {
	return server.Config{
		Address: appCfg.Listen.BindIP,
		Port:    appCfg.Listen.Port}
}

func GetKafkaWriter(cfg config.KafkaConfig) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return w
}
