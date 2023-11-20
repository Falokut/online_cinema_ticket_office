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

// CreateAccountAndProfile creates a new account and profile in the database.
func (r *postgreRepository) CreateAccountAndProfile(ctx context.Context, account model.CreateAccountAndProfile) error {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"PostgreRepository.CreateAccountAndProfile")
	span.SetTag("database", "postgre")
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
	res, _ := tx.Exec(query, id, account.Email, account.Username, account.RegistrationDate)
	num, err := res.RowsAffected()
	if num == 0 && err != nil {
		return errors.New("no rows affected")
	}

	return tx.Commit()
}

// IsAccountWithEmailExist checks if an account with the given email exists in the database.
// It returns a boolean indicating the existence and an error, if any.
func (r *postgreRepository) IsAccountWithEmailExist(ctx context.Context, email string) (bool, error) {
	// Start a new span for tracing.
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.IsAccountWithEmailExist")
	span.SetTag("database", "postgre")
	defer span.Finish() // Finish the span when the function ends.

	// Prepare the SQL query to check for the existence of the account.
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var UUID string
	// Execute the query to check for the existence of the account with the given email.
	err := r.db.Get(&UUID, query, email)
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
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.GetUserByEmail")
	span.SetTag("database", "postgre")
	defer span.Finish()

	// Prepare the SQL query to retrieve the user account based on the provided email.
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 LIMIT 1;", accountTableName)

	var acc model.Account
	// Execute the query to retrieve the user account.
	err := r.db.Get(&acc, query, email)
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
	span.SetTag("database", "postgre")
	defer span.Finish() // Finish the span when the function ends.

	// Prepare the SQL query to update the password hash of the account with the given email.
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE email=$2;", accountTableName)

	// Execute the query to update the password hash.
	res, err := r.db.Exec(query, passwordHash, email)
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
func (r *postgreRepository) DeleteAccount(ctx context.Context, AccountID string) error {
	// Start a new span for tracing.
	span, _ := opentracing.StartSpanFromContext(ctx, "PostgreRepository.DeleteAccount")
	span.SetTag("database", "postgre")
	defer span.Finish() // Finish the span when the function ends.

	// Prepare the SQL query to delete the account with the given ID.
	query := fmt.Sprintf("DELETE * FROM %s WHERE id=$1;", accountTableName)

	// Execute the query to delete the account.
	_, err := r.db.Exec(query, AccountID)

	return err // Return any error that occurs during the query execution.
}
