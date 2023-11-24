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

// NewPostgreDB creates a new connection to the PostgreSQL database.
func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	return db, err
}

// NewAccountRepository creates a new instance of the postgreRepository using the provided database connection.
func NewAccountRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

// Shutdown closes the database connection.
func (r *postgreRepository) Shutdown() error {
	return r.db.Close()
}

// CreateAccount creates a new account in the database.
func (r *postgreRepository) CreateAccount(ctx context.Context, account model.Account) (*sql.Tx, string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.CreateAccountAndProfile")

	defer span.Finish()

	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, registration_date) VALUES ($1, $2, $3) RETURNING id;", accountTableName)
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, "", err
	}

	row := tx.QueryRowContext(ctx, query, account.Email, account.Password, account.RegistrationDate)

	var id string
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return nil, "", err
	}

	return tx, id, nil
}

// IsAccountWithEmailExist checks if an account with the given email exists in the database.
// It returns a boolean indicating the existence and an error, if any.
func (r *postgreRepository) IsAccountWithEmailExist(ctx context.Context, email string) (bool, error) {
	// Start a new span for tracing.
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.IsAccountWithEmailExist")

	defer span.Finish() // Finish the span when the function ends.

	// Prepare the SQL query to check for the existence of the account.
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var UUID string
	// Execute the query to check for the existence of the account with the given email.
	err := r.db.GetContext(ctx, &UUID, query, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err // Return false and the error if an error other than sql.ErrNoRows occurs.
	}
	if err == sql.ErrNoRows {
		return false, nil // Return false if no rows were found (account does not exist).
	}

	return true, nil // If no error occurred, the account exists.
}

// GetAccountByEmail retrieves a account from the database based on the provided email.
// It returns the retrieved account and an error, if any.
func (r *postgreRepository) GetAccountByEmail(ctx context.Context, email string) (model.Account, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.GetAccountByEmail")
	defer span.Finish()

	// Prepare the SQL query to retrieve the user account based on the provided email.
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var acc model.Account
	// Execute the query to retrieve the user account.
	err := r.db.GetContext(ctx, &acc, query, email)
	if err == sql.ErrNoRows {
		return model.Account{}, errors.New("user with this email doesn't exist") // Return an error if no account was found with the provided email.
	}

	return acc, err // Return the retrieved account and any error that occurred during retrieval.
}

// ChangePassword updates the password hash of an account with the given email in the database.
// It takes the email and the new password hash as input and returns an error, if any.
func (r *postgreRepository) ChangePassword(ctx context.Context, email string, passwordHash string) error {
	// Start a new span for tracing.
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.ChangePassword")

	defer span.Finish() // Finish the span when the function ends.

	// Prepare the SQL query to update the password hash of the account with the given email.
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE email=$2;", accountTableName)

	// Execute the query to update the password hash.
	res, err := r.db.ExecContext(ctx, query, passwordHash, email)
	if err != nil {
		return err // Return the error if the query execution fails.
	}

	// Get the number of rows affected by the update.
	num, err := res.RowsAffected()
	if err != nil || num == 0 {
		return errors.New("rows are not affected") // Return an error if no rows are affected by the update.
	}

	return nil // Return nil to indicate success.
}

// DeleteAccount deletes the account with the given ID from the database.
// It takes the ID of the account as input and returns an error, if any.
func (r *postgreRepository) DeleteAccount(ctx context.Context, AccountID string) (*sql.Tx, error) {
	// Start a new span for tracing.
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.DeleteAccount")
	defer span.Finish() // Finish the span when the function ends.

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	// Prepare the SQL query to delete the account with the given ID.
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", accountTableName)

	// Execute the query to delete the account.
	_, err = tx.ExecContext(ctx, query, AccountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return tx, nil
}
