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
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port"`
	Listen          struct {
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
		StorageAddr            string                         `yaml:"storage_addr" env:"IMAGE_STORAGE_ADDRESS"`
		ImageResizeType        images_resizer.ImageResizeType `yaml:"image_resize_type" env:"IMAGE RESIZE_TYPE"`
		ImageResizeMethod      string                         `yaml:"image_resize_method" env:"IMAGE_RESIZE_METHOD"`
		BaseProfilePictureUrl  string                         `yaml:"base_profile_picture_url" env:"BASE_PROFILE_PICTURE_URL"`
		ProfilePictureCategory string                         `yaml:"profile_picture_category" env:"PROFILE_PICTURE_CATEGORY"`

		ImageWidth     uint `yaml:"image_width" env:"IMAGE_WIDTH"`
		ImageHeight    uint `yaml:"image_height" env:"IMAGE_HEIGHT"`
		MaxImageWidth  uint `yaml:"max_image_width" env:"MAX_IMAGE_WIDTH"`
		MaxImageHeight uint `yaml:"max_image_height" env:"MAX_IMAGE_HEIGHT"`
		MinImageWidth  uint `yaml:"min_image_width" env:"MIN_IMAGE_WIDTH"`
		MinImageHeight uint `yaml:"min_image_height" env:"MIN_IMAGE_HEIGHT"`
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
