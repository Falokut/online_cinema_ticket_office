package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

const (
	accountTableName  = "accounts"
	profilesTableName = "profiles"
)

type postgreRepository struct {
	db *sqlx.DB
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	return db, err
}

func NewAccountRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func (r *postgreRepository) ShutDown() error {
	return r.db.Close()
}

func (r *postgreRepository) CreateAccountAndProfile(ctx context.Context, account model.CreateAccountAndProfile) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.CreateAccountAndProfile")
	span.SetTag("custom-tag", "database")
	defer span.Finish()

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, registration_date) VALUES ($1, $2, $3) RETURNING id;", accountTableName)

	row := tx.QueryRow(query, account.Email, account.Password, account.RegistrationDate)

	var id string
	if err := row.Scan(&id); err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (account_id, email, username, registration_date) VALUES ($1, $2, $3, $4);", profilesTableName)
	res, err := tx.Exec(query, id, account.Email, account.Username, account.RegistrationDate)
	num, err := res.RowsAffected()
	if num == 0 && err != nil {
		return errors.New("No rows affected")
	}

	return tx.Commit()
}

func (r *postgreRepository) IsAccountWithEmailExist(ctx context.Context, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.IsAccountWithEmailExist")
	span.SetTag("custom-tag", "database")
	defer span.Finish()

	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var UUID string
	err := r.db.Get(&UUID, query, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (r *postgreRepository) GetUserByEmail(ctx context.Context, email string) (model.Account, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostgreRepository.GetUserByEmail")
	span.SetTag("custom-tag", "database")
	defer span.Finish()
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var acc model.Account
	err := r.db.Get(&acc, query, email)
	if err == sql.ErrNoRows {
		return model.Account{}, errors.New("User with this email doesn't exist.")
	}

	return acc, err
}

func (r *postgreRepository) ChangePassword(ctx context.Context, email string, password_hash string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostgreRepository.ChangePassword")
	span.SetTag("custom-tag", "database")
	defer span.Finish()
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE email=$2;", accountTableName)

	res, err := r.db.Exec(query, password_hash, email)
	if err != nil {
		return err
	}

	num, err := res.RowsAffected()
	if err != nil || num == 0 {
		return errors.New("Rows are not affected")
	}

	return nil
}

func (r *postgreRepository) DeleteAccount(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PostgreRepository.DeleteAccount")
	span.SetTag("custom-tag", "database")
	defer span.Finish()

	query := fmt.Sprintf("DELETE * FROM %s WHERE id=$1;", accountTableName)
	_, err := r.db.Exec(query, id)

	return err
}
