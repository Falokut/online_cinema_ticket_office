package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/email_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/email_service/internal/email"
	"github.com/Falokut/online_cinema_ticket_office/email_service/pkg/logging"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.GetLogger()
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	appCfg := config.GetConfig()
	log_level, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Logger.SetLevel(log_level)

	mailSender := email.NewMailSender(appCfg.MailSenderCfg, logger)

	logger.Infoln("kafka consumer initializing")
	kafkaReader := NewKafkaReader(*appCfg)

	logger.Infoln("worker initializing")
	mailWorker := email.NewMailWorker(mailSender, logger, appCfg.MailWorkerCfg, kafkaReader, appCfg.MaxWorkersCount)
	go func() {
		mailWorker.Run()
	}()

	var wg sync.WaitGroup
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	kafkaReader.Close()
	wg.Wait()
	logger.Infoln("Shutted down successfully")
}

func NewKafkaReader(appCfg config.Config) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        appCfg.KafkaConfig.Brokers,
		GroupID:        appCfg.KafkaConfig.GroupID,
		Topic:          appCfg.KafkaConfig.Topic,
		MaxBytes:       appCfg.KafkaConfig.MaxBytes,
		Logger:         logging.GetLogger(),
		MaxAttempts:    4,
		MaxWait:        time.Minute,
		ReadBackoffMax: time.Millisecond * 100,
		StartOffset:    kafka.LastOffset,
		QueueCapacity:  appCfg.KafkaConfig.QueueCapacity,
	})
	return r
}
