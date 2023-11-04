package repository

import (
	"context"
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
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

func (r *postgreRepository) GetUserProfile(ctx context.Context, AccountID string) (model.UserProfile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.GetUserProfile")
	defer span.Finish()

	query := fmt.Sprintf("SELECT username, email, profile_picture_id, registration_date FROM %s WHERE account_id=$1 LIMIT 1;",
		profilesTableName)
	var Profile model.UserProfile
	err := r.db.Get(&Profile, query, AccountID)

	return Profile, err
}

func (r *postgreRepository) UpdateProfilePictureID(ctx context.Context, AccountID string, PictureID string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.GetUserProfile")
	defer span.Finish()

	query := fmt.Sprintf("UPDATE %s SET profile_picture_id=$1 WHERE account_id=$2;",
		profilesTableName)

	_, err := r.db.Exec(query, PictureID, AccountID)
	return err
}

func (r *postgreRepository) Shutdown() error {
	return r.db.Close()
}
