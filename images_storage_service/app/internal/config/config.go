package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/jaeger"

	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel             string `yaml:"log_level" env:"LOG_LEVEL"`
	BaseLocalStoragePath string `yaml:"base_local_storage_path" env:"BASE_LOCAL_STORAGE_PATH"`

	// max image size in bytes
	MaxImageSize int `yaml:"max_image_size" env:"MAX_IMAGE_SIZE"`
	Listen       struct {
		Host           string   `yaml:"host" env:"HOST"`
		Port           string   `yaml:"port" env:"PORT"`
		AllowedHeaders []string `yaml:"allowed_headers"`
		Mode           string   `yaml:"server_mode" env:"SERVER_MODE"`
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`
	JaegerConfig jaeger.Config `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
