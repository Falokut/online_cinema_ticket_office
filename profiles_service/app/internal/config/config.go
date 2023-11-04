package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
	Listen   struct {
		BindIP string `yaml:"bind_ip" env:"BIND_IP"`
		Port   uint16 `yaml:"port" env:"LISTEN_PORT"`
	} `yaml:"listen"`
	PrometheusConfig struct {
		Address string `yaml:"address" env:"PROMETHEUS_ADDRESS"`
		Name    string `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
	} `yaml:"prometheus"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
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

		if err := cleanenv.ReadConfig(configsPath+".env", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})
	return instance
}
