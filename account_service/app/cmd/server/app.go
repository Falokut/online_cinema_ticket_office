package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/service"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
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

	logger.Logger.SetLevel(log_level)
	logger.Infoln("Registration cache initializing")
	regCacheOpt := appCfg.RegistrationCacheOptions.ConvertToRedisOptions()
	registrationCache, err := repository.NewRedisRegistrationCache(regCacheOpt, logger)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis registration cache is not established: %s options: %v", err.Error(), regCacheOpt)
		return
	}

	logger.Infoln("Token cache initializing")
	sessionCacheOpt := appCfg.SessionCacheOptions.ConvertToRedisOptions()
	accountSessionCacheOpt := appCfg.AccountSessionsCacheOptions.ConvertToRedisOptions()
	sessionCache, err := repository.NewSessionCache(sessionCacheOpt, accountSessionCacheOpt, logger, appCfg.SessionsTTL)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the redis token cache is not established: %s options: %v", err.Error(), sessionCacheOpt)
		return
	}

	logger.Infoln("database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}

	redisRepo := repository.NewCacheRepository(registrationCache, sessionCache)

	logger.Infoln("repository initializing")
	repo := repository.NewAccountRepository(database)

	logger.Infoln("service initializing")
	service := service.NewAccountService(repo, logger, redisRepo, GetKafkaWriter(appCfg.EmailKafka), appCfg)

	s := server.NewServer(logger, service)
	logger.Infoln("server initializing")
	go func() {
		logger.Infoln("server running")
		if err := s.Run(getServerConfig(appCfg)); err != nil {
			log.Fatalf("%s", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.ShutDown()
	logger.Infoln("Shutted down successfully")
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
