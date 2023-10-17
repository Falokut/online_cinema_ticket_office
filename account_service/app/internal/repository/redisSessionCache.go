package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/redis/go-redis/v9"
)

type redisSessionCache struct {
	sessions_rdb         *redis.Client
	account_sessions_rdb *redis.Client
	logger               logging.Logger
	SessionTTL           time.Duration
}

func NewSessionCache(sessionCacheOpt *redis.Options, accountSessionsOpt *redis.Options, logger logging.Logger, SessionTTL time.Duration) (*redisSessionCache, error) {
	logger.Println("Creating session cache client")
	sessions_rdb := redis.NewClient(sessionCacheOpt)
	if sessions_rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Println("Pinging seccion cache client")
	_, err := sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	logger.Println("Creating account sessions cache client")
	account_sessions_rdb := redis.NewClient(accountSessionsOpt)
	if account_sessions_rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Println("Pinging seccion cache client")
	_, err = account_sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	return &redisSessionCache{sessions_rdb: sessions_rdb, account_sessions_rdb: account_sessions_rdb, logger: logger, SessionTTL: SessionTTL}, nil
}

func (r *redisSessionCache) ShutDown() {
	r.logger.Println("Token cache repository shutting down")
	r.sessions_rdb.Close()
}

func (r *redisSessionCache) CacheSession(toCache model.SessionCache) error {
	r.logger.Debug("Marshalling data")
	serialized, err := json.Marshal(toCache)
	if err != nil {
		return err
	}

	_, err = r.sessions_rdb.Set(context.Background(), toCache.SessionID, serialized, r.SessionTTL).Result()
	if err != nil {
		return errors.New("Can't cache session data")
	}

	if err := r.cacheAccountSession(toCache); err != nil {
		r.sessions_rdb.Del(context.Background(), toCache.SessionID)
		return err
	}

	return nil
}

func (r *redisSessionCache) TerminateSessions(sessionsID []string, accountID string) error {
	AccountSessions, err := r.GetSessionsForAccount(accountID)
	if err != nil {
		return err
	}

	if err := r.sessions_rdb.Del(context.Background(), sessionsID...).Err(); err != nil {
		return err
	}

	for _, sessionID := range sessionsID {
		delete(AccountSessions, sessionID)
	}

	if err := r.UpdateSessionsForAccount(AccountSessionsCache{AccountSessions}, accountID); err != nil {
		return err
	}

	return nil
}

func (r *redisSessionCache) cacheAccountSession(toCache model.SessionCache) error {
	body, err := r.account_sessions_rdb.Get(context.Background(), toCache.AccountID).Bytes()
	if err != nil && err != redis.Nil {
		r.logger.Debugf("Can't read the bytes, %s", err.Error())
		return err
	}

	var sessionsCache AccountSessionsCache
	sessionsCache.Sessions = make(map[string]sessionInfo)

	if err != redis.Nil {
		r.logger.Debug("Unmarshal data")
		if err := json.Unmarshal(body, &sessionsCache); err != nil {
			return err
		}
	}
	sessionsCache.Sessions[toCache.SessionID] = sessionInfo{ClientIP: toCache.ClientIP,
		SessionInfo: toCache.SessionInfo, LastActivity: toCache.LastActivity}

	return r.UpdateSessionsForAccount(sessionsCache, toCache.AccountID)
}

func (r *redisSessionCache) GetSessionCache(sessionID string) (model.SessionCache, error) {
	body, err := r.sessions_rdb.Get(context.Background(), sessionID).Bytes()
	if err == redis.Nil {
		return model.SessionCache{}, errors.New("Session not found")
	}

	r.logger.Debug("Unmarshal cache data")
	var session model.SessionCache
	if err := json.Unmarshal(body, &session); err != nil {
		return model.SessionCache{}, err
	}

	return session, nil
}

func (r *redisSessionCache) UpdateLastActivityForSession(sessionID string, activityTime time.Time) error {
	cache, err := r.GetSessionCache(sessionID)
	if err != nil {
		return err
	}
	cache.LastActivity = activityTime
	return r.CacheSession(cache)
}

type sessionInfo struct {
	ClientIP     string    `json:"client_ip"`
	SessionInfo  string    `json:"session_info"` // like device or browser
	LastActivity time.Time `json:"last_activity"`
}

type AccountSessionsCache struct {
	Sessions map[string]sessionInfo `json:"sessions"`
}

func (r *redisSessionCache) GetSessionsForAccount(accountID string) (map[string]sessionInfo, error) {
	body, err := r.account_sessions_rdb.Get(context.Background(), accountID).Bytes()
	if err == redis.Nil {
		return map[string]sessionInfo{}, nil // not error, just empty sessions map for account
	}

	r.logger.Debug("Unmarshal cache data")
	var sessions AccountSessionsCache
	if err := json.Unmarshal(body, &sessions); err != nil {
		return map[string]sessionInfo{}, err
	}

	return sessions.Sessions, nil
}

func (r *redisSessionCache) UpdateSessionsForAccount(sessions AccountSessionsCache, accountID string) error {
	if len(sessions.Sessions) == 0 {
		err := r.account_sessions_rdb.Del(context.Background(), accountID).Err()
		if err == redis.Nil {
			return nil
		}
		return err
	}

	r.logger.Debug("Marshalling data")
	serialized, err := json.Marshal(sessions)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	r.logger.Debug("Caching data")
	_, err = r.account_sessions_rdb.Set(context.Background(), accountID, serialized, r.SessionTTL).Result()
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}
