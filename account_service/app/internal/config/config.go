package config

import (
	"sync"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
)

type token struct {
	TTL    time.Duration `yaml:"TTL"`
	Secret string        `yaml:"secret"`
}

type redisOptions struct {
	Network  string `yaml:"network"`
	Addr     string `yaml:"addr"`
	Password string `env:"PASSWORD"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func (r redisOptions) ConvertToRedisOptions() *redis.Options {
	return &redis.Options{
		Network:  r.Network,
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	}
}

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
	Listen   struct {
		BindIP string `yaml:"bind_ip" env:"BIND_IP"`
		Port   uint16 `yaml:"port" env:"PORT"`
	}
	SessionsTTL time.Duration       `yaml:"sessions_ttl"` // The lifetime of an inactive session in the cache
	DBConfig    repository.DBConfig `yaml:"db_config"`
	EmailKafka  KafkaConfig         `yaml:"email_kafka_config"`

	RegistrationCacheOptions    redisOptions `yaml:"redis_registration_options"`
	SessionCacheOptions         redisOptions `yaml:"session_cache_options"`
	AccountSessionsCacheOptions redisOptions `yaml:"account_sessions_cache_options"`

	Crypto struct {
		BcryptCost int `yaml:"bcrypt_cost" enb:"BCRYPT_COST"`
	} `yaml:"crypto"`
	JWT struct {
		ChangePasswordToken token `yaml:"change_password_token"`
		VerifyAccountToken  token `yaml:"verify_account_token"`
	} `yaml:"JWT"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		err := cleanenv.ReadEnv(instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}

		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}

		if err := cleanenv.ReadConfig("configs/secrets.env.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
