package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/server"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/service"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/logging"
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

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	logger.Println("database initializing")
	database, err := repository.NewPostgreDB(getRepositoryConfig(appCfg))

	if err != nil {
		logger.Fatalf("shuting down, connection to the database is not established: %s", err.Error())
		return
	}

	logger.Println("repository initializing")
	repo := repository.NewProfileRepository(database)

	service := service.NewProfileService(repo, logger)
	s := server.Server{}

	logger.Println("server initializing")
	go func() {
		if err := s.Run(getServerConfig(appCfg), service); err != nil {
			log.Fatalf("%s", err.Error())
		}
	}()

	var wg sync.WaitGroup
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
	if err := repo.Shutdown(); err != nil {
		logger.Error(err)
	}

	wg.Wait()
	logger.Infoln("Shutted down successfully")
}

func getRepositoryConfig(appCfg *config.Config) repository.Config {
	return repository.Config{
		Host:     appCfg.DBConfig.Host,
		Port:     appCfg.DBConfig.Port,
		Username: appCfg.DBConfig.Username,
		Password: appCfg.DBConfig.Password,
		DBName:   appCfg.DBConfig.DBName,
		SSLMode:  appCfg.DBConfig.SSLMode}
}

func getServerConfig(appCfg *config.Config) server.Config {
	return server.Config{
		Address: appCfg.Listen.BindIP,
		Port:    appCfg.Listen.Port}
}
