package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	EnableTLS bool   `yaml:"enable_TLS"`
	LogLevel  string `yaml:"log_level"`

	Listen struct {
		BindIP string `yaml:"bind_ip"`
		Port   uint16 `yaml:"port"`
	}
	DBConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `env:"DB_PASSWORD,env-required"  env-default:"password" `
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
	}
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
	})
	return instance
}
