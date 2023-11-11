package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/service"
	jaegerTracer "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/metrics"
	"github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
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

	logger.Info("Registration cache initializing")
	regCacheOpt := appCfg.RegistrationCacheOptions.ConvertToRedisOptions()
	registrationCache, err := repository.NewRedisRegistrationCache(regCacheOpt, logger)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis registration cache is not established: %s options: %v",
			err.Error(), regCacheOpt)
	}
	defer registrationCache.ShutDown()

	logger.Info("Sessions cache initializing")
	sessionCacheOpt := appCfg.SessionCacheOptions.ConvertToRedisOptions()
	accountSessionCacheOpt := appCfg.AccountSessionsCacheOptions.ConvertToRedisOptions()
	sessionsCache, err := repository.NewSessionCache(sessionCacheOpt,
		accountSessionCacheOpt, logger, appCfg.SessionsTTL)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis token cache is not established: %s options: %v",
			err.Error(), sessionCacheOpt)
	}
	defer sessionsCache.ShutDown()

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}
	defer database.Close()

	redisRepo := repository.NewCacheRepository(registrationCache, sessionsCache)

	logger.Info("Repository initializing")
	repo := repository.NewAccountRepository(database)

	logger.Info("Service initializing")
	service := service.NewAccountService(repo,
		logger, redisRepo, getKafkaWriter(appCfg.EmailKafka), appCfg, metric)

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
		Mode:           cfg.Listen.Mode,
		Host:           cfg.Listen.Host,
		Port:           cfg.Listen.Port,
		AllowedHeaders: cfg.Listen.AllowedHeaders,
	}
}
func getKafkaWriter(cfg config.KafkaConfig) *kafka.Writer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return w
}
