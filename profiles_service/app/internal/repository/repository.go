package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/model"
	"github.com/jmoiron/sqlx"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type ProfileRepository interface {
	CreateUserProfile(ctx context.Context, profile model.UserProfile) error
	DeleteUserProfile(ctx context.Context, AccountID string) error
	GetUserProfile(ctx context.Context, AccountID string) (model.UserProfile, error)
	GetProfilePictureID(ctx context.Context, AccountID string) (string, error)
	UpdateProfilePictureID(ctx context.Context, AccountID string, PictureID string) error
	GetEmail(ctx context.Context, AccountID string) (string, error)
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}
