package repository

import (
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type postgreRepository struct {
	db *sqlx.DB
}

const (
	profilesTableName = "profiles"
)

func NewProfileRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func (r *postgreRepository) GetUserProfile(UUID string) (model.UserProfile, error) {
	query := fmt.Sprintf("SELECT username, email, profile_picture_id, registration_date FROM %s WHERE id=$1 LIMIT 1;", profilesTableName)
	var Profile model.UserProfile
	err := r.db.Get(&Profile, query, UUID)

	if err != nil {
		return model.UserProfile{}, err
	}

	return Profile, nil
}

func (r *postgreRepository) UpdateProfilePicture(UUID string, PictureID string) error {
	query := fmt.Sprintf("UPDATE %s SET profile_picture_id=$1 WHERE id=$2 RETURNING profile_picture_id;", profilesTableName)

	_, err := r.db.Exec(query, PictureID, UUID)
	return err
}

func (r *postgreRepository) Shutdown() error {
	return r.db.Close()
}
