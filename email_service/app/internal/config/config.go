package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/email_service/internal/email"
	"github.com/Falokut/online_cinema_ticket_office/email_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	EnableTLS     bool                   `yaml:"enable_TLS" env:"ENABLE_TLS"`
	LogLevel      string                 `yaml:"log_level" env:"LOG_LEVEL"`
	MailSenderCfg email.MailSenderConfig `yaml:"mail_sender"`
	MailWorkerCfg email.MailWorkerConfig `yaml:"mail_worker"`
	KafkaConfig   struct {
		Brokers       []string `yaml:"brokers"`
		GroupID       string   `yaml:"group_id"`
		Topic         string   `yaml:"topic"`
		MaxBytes      int      `yaml:"max_bytes"`
		QueueCapacity int      `yaml:"queque_capacity"`
	} `yaml:"kafka_config"`
}

const configsPath string = "configs/"

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		logger := logging.GetLogger()
		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})
	return instance
}
