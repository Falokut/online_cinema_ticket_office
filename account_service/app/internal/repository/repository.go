package repository

import (
	"context"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/model"
)

type AccountRepository interface {
	CreateAccountAndProfile(ctx context.Context, account model.CreateAccountAndProfile) error
	IsAccountWithEmailExist(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (model.Account, error)
	ChangePassword(ctx context.Context, email string, password_hash string) error
	DeleteAccount(ctx context.Context, id string) error
	ShutDown() error
}

type CachedAccount struct {
	Username string
	Password string
}

type RegistrationCacheRepository interface {
	IsAccountInCache(ctx context.Context, email string) (bool, error)
	CacheAccount(ctx context.Context, email string, Account CachedAccount, NonActivatedAccountTTL time.Duration) error
	DeleteAccountFromCache(ctx context.Context, email string) error
	GetCachedAccount(ctx context.Context, email string) (CachedAccount, error)
	ShutDown() error
}

type SeccionsCacheRepository interface {
	CacheSession(ctx context.Context, toCache model.SessionCache) error
	TerminateSessions(ctx context.Context, sessionsID []string, accountID string) error
	UpdateLastActivityForSession(ctx context.Context, cachedSession model.SessionCache, sessionID string, LastActivityTime time.Time) error
	GetSessionCache(ctx context.Context, sessionID string) (model.SessionCache, error)
	GetSessionsForAccount(ctx context.Context, accountID string) (map[string]sessionInfo, error)
	ShutDown() error
}

type CacheRepo struct {
	RegistrationCache RegistrationCacheRepository
	SessionsCache     SeccionsCacheRepository
}

func NewCacheRepository(account RegistrationCacheRepository, SessionsCache SeccionsCacheRepository) CacheRepo {
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
