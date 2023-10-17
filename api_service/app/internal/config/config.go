package config

import (
	"sync"

	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Mode string `yaml:"mode" env-required:"true" env-default:"debug"`
	JWT  struct {
		Secret          string `env:"JWT_SIGN_SECRET"`
		TokenTTL        int    `yaml:"tokenTTL"`
		RefreshTokenTTL int    `yaml:"refresh_tokenTTL"`
	}

	Listen struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
		Domain string `yaml:"domain"`
	}

	Contexts struct {
		UserContext string `yaml:"user_context"`
	} `yaml:"contexts"`

	ImageStorage struct {
		URL string `yaml:"url"`
	} `yaml:"image_storage_service"`
	MoviePostersService struct {
		URL string `yaml:"url"`
	} `yaml:"movie_posters_service"`

	UserService struct {
		URL string `yaml:"url"`
	} `yaml:"user_service"`
	AccountService struct {
		URL string `yaml:"url"`
	} `yaml:"account_service"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}

		logger.Info("read application env config")
		if err := cleanenv.ReadConfig("configs/.env", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
