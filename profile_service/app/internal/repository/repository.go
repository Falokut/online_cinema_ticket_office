package repository

import (
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository interface {
	GetUserProfile(UUID string) (model.UserProfile, error)
	UpdateProfilePicture(UUID string, PictureID string) error
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgreDB(cfg Config) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		logging.GetLogger().Error(conStr)
		return nil, err
	}

	return db, nil
}
