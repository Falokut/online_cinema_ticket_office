package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/model"
)

// AccountRepository provides methods to interact with user accounts in the database.
type AccountRepository interface {
	// CreateAccount creates a new account in the database.
	CreateAccount(ctx context.Context, account model.Account) (*sql.Tx, string, error)

	// IsAccountWithEmailExist checks if an account with the given email exists in the database.
	IsAccountWithEmailExist(ctx context.Context, email string) (bool, error)

	// GetAccountByEmail retrieves an account from the database using the email.
	GetAccountByEmail(ctx context.Context, email string) (model.Account, error)

	// ChangePassword updates the password hash of an account with the given email in the database.
	ChangePassword(ctx context.Context, email string, passwordHash string) error

	// DeleteAccount removes the account with the given ID from the database.
	DeleteAccount(ctx context.Context, id string) (*sql.Tx, error)

	// Shutdown performs cleanup and shuts down the repository.
	Shutdown() error
}

// CachedAccount represents the cached account data.
type CachedAccount struct {
	Username string
	Password string
}

// RegistrationCacheRepository provides methods to interact with the registration cache.
type RegistrationCacheRepository interface {
	// IsAccountInCache checks if the account associated with the given email is present in the cache.
	// It returns true if the account is found in the cache, otherwise false.
	// An error is returned if there is an issue while checking the cache.
	IsAccountInCache(ctx context.Context, email string) (bool, error)

	// CacheAccount caches the account with the given email and its details with a specified time-to-live duration.
	CacheAccount(ctx context.Context, email string, Account CachedAccount, NonActivatedAccountTTL time.Duration) error

	// DeleteAccountFromCache removes the account with the given email from the cache.
	DeleteAccountFromCache(ctx context.Context, email string) error

	// GetCachedAccount retrieves the cached account data using the email.
	GetCachedAccount(ctx context.Context, email string) (CachedAccount, error)

	PingContext(ctx context.Context) error

	// Shutdown performs cleanup and shuts down the registration cache repository.
	Shutdown() error
}

var (
	ErrSessionNotFound = errors.New("session not found")
)

// SessionsCacheRepository provides methods to interact with the sessions cache.
type SessionsCacheRepository interface {
	// CacheSession caches the session data.
	CacheSession(ctx context.Context, toCache model.SessionCache) error

	// TerminateSessions terminates the specified sessions for the given account ID.
	TerminateSessions(ctx context.Context, sessionsID []string, accountID string) error

	// UpdateLastActivityForSession updates the last activity time for the session with the given ID.
	UpdateLastActivityForSession(ctx context.Context, cachedSession model.SessionCache, sessionID string, LastActivityTime time.Time) error

	// GetSessionCache retrieves the cached session data for a specific session ID.
	GetSessionCache(ctx context.Context, sessionID string) (model.SessionCache, error)

	// GetSessionsForAccount fetches the sessions associated with a particular account ID.
	GetSessionsForAccount(ctx context.Context, accountID string) (map[string]SessionInfo, error)

	// GetSessionsForAccount fetches the sessions associated with a particular account ID.
	GetSessionsList(ctx context.Context, accountID string) ([]string, error)

	PingContext(ctx context.Context) error

	// Shutdown stops the cache operations.
	Shutdown() error
}

type CacheRepo struct {
	RegistrationCache RegistrationCacheRepository
	SessionsCache     SessionsCacheRepository
}

func NewCacheRepository(account RegistrationCacheRepository, SessionsCache SessionsCacheRepository) CacheRepo {
	return CacheRepo{RegistrationCache: account, SessionsCache: SessionsCache}
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD,env-required" env-default:"password"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}
