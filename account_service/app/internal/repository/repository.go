package repository

import (
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/model"
)

type AccountRepository interface {
	CreateAccountAndProfile(account model.CreateAccountAndProfile) error
	IsAccountWithEmailExist(email string) (bool, error)
	GetUserByEmail(email string) (model.Account, error)
	ChangePassword(email string, password_hash string) error
	DeleteAccount(id string) error
	ShutDown()
}

type CachedAccount struct {
	Username string
	Password string
}

type RegistrationCacheRepository interface {
	IsAccountInCache(email string) (bool, error)
	CacheAccount(email string, Account CachedAccount, NonActivatedAccountTTL time.Duration) error
	DeleteAccountFromCache(email string) error
	GetCachedAccount(email string) (CachedAccount, error)
	ShutDown()
}

type SeccionCacheRepository interface {
	CacheSession(toCache model.SessionCache) error
	TerminateSessions(sessionsID []string, accountID string) error
	UpdateLastActivityForSession(sessionID string, activityTime time.Time) error
	GetSessionCache(sessionID string) (model.SessionCache, error)
	GetSessionsForAccount(accountID string) (map[string]sessionInfo, error)
	ShutDown()
}

type CacheRepo struct {
	RegistrationCache RegistrationCacheRepository
	SessionCache      SeccionCacheRepository
}

func NewCacheRepository(account RegistrationCacheRepository, sessionCache SeccionCacheRepository) CacheRepo {
	return CacheRepo{RegistrationCache: account, SessionCache: sessionCache}
}

func (r *CacheRepo) ShutDown() {
	r.RegistrationCache.ShutDown()
	r.SessionCache.ShutDown()
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password" env:"DB_PASSWORD,env-required"  env-default:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}
