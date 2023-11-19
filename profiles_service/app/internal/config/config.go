package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_resizer"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/jaeger"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
	Listen   struct {
		Host           string   `yaml:"host" env:"HOST"`
		Port           string   `yaml:"port" env:"PORT"`
		Mode           string   `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
		AllowedHeaders []string `yaml:"allowed_headers"`               // Need for REST API gateway, list of metadata headers
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	ImageService struct {
		StorageAddr       string                         `yaml:"storage_addr"`
		ImageWidth        uint                           `yaml:"image_width"`
		ImageHeight       uint                           `yaml:"image_height"`
		ImageResizeType   images_resizer.ImageResizeType `yaml:"image_resize_type"`
		ImageResizeMethod string                         `yaml:"image_resize_method"`

		BaseProfilePictureUrl  string `yaml:"base_profile_picture_url"`
		ProfilePictureCategory string `yaml:"profile_picture_category"`
	} `yaml:"image_service"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

const configsPath string = "configs/"

var instance *Config
var once sync.Once

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
