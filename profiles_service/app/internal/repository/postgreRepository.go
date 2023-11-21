package repository

import (
	"context"
	"database/sql"
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

func (r *postgreRepository) GetProfilePictureID(ctx context.Context, AccountID string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.GetProfilePictureID")
	defer span.Finish()
	query := fmt.Sprintf("SELECT profile_picture_id FROM %s WHERE account_id=$1 LIMIT 1;",
		profilesTableName)

	var PictureID []sql.NullString
	err := r.db.Select(&PictureID, query, AccountID)
	if err != nil {
		return "", err
	}

	if len(PictureID) == 0 || !PictureID[0].Valid {
		return "", nil
	}

	return PictureID[0].String, nil
}
func (r *postgreRepository) GetEmail(ctx context.Context, AccountID string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.GetEmail")
	defer span.Finish()
	query := fmt.Sprintf("SELECT email FROM %s WHERE account_id=$1 LIMIT 1;",
		profilesTableName)

	var Email []sql.NullString
	err := r.db.Select(&Email, query, AccountID)
	if err != nil {
		return "", err
	}

	if len(Email) == 0 || !Email[0].Valid {
		return "", nil
	}

	return Email[0].String, nil
}

func (r *postgreRepository) Shutdown() error {
	return r.db.Close()
}
