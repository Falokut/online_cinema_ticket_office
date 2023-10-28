package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type redisSessionsCache struct {
	sessions_rdb         *redis.Client
	account_sessions_rdb *redis.Client
	logger               logging.Logger
	SessionTTL           time.Duration
}

func NewSessionCache(sessionCacheOpt *redis.Options, accountSessionsOpt *redis.Options, logger logging.Logger, SessionTTL time.Duration) (*redisSessionsCache, error) {
	logger.Infoln("Creating session cache client")
	sessions_rdb := redis.NewClient(sessionCacheOpt)
	if sessions_rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Infoln("Pinging seccion cache client")
	_, err := sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	logger.Infoln("Creating account sessions cache client")
	account_sessions_rdb := redis.NewClient(accountSessionsOpt)
	if account_sessions_rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Infoln("Pinging seccion cache client")
	_, err = account_sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	return &redisSessionsCache{sessions_rdb: sessions_rdb, account_sessions_rdb: account_sessions_rdb, logger: logger, SessionTTL: SessionTTL}, nil
}

func (r *redisSessionsCache) ShutDown() {
	r.logger.Infoln("Token cache repository shutting down")
	r.sessions_rdb.Close()
}

func (r *redisSessionsCache) CacheSession(ctx context.Context, toCache model.SessionCache) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.CacheSession")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	r.logger.Debug("Marshalling data")
	serialized, err := json.Marshal(toCache)
	if err != nil {
		return err
	}

	r.logger.Debug("Caching sessions data")
	_, err = r.sessions_rdb.Set(ctx, toCache.SessionID, serialized, r.SessionTTL).Result()
	if err != nil {
		return errors.New("Can't cache session data")
	}

	if err := r.cacheAccountSession(ctx, toCache); err != nil {
		r.sessions_rdb.Del(ctx, toCache.SessionID)
		return err
	}

	return nil
}

func (r *redisSessionsCache) TerminateSessions(ctx context.Context, sessionsID []string, accountID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.TerminateSessions")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	AccountSessions, err := r.GetSessionsForAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := r.sessions_rdb.Del(ctx, sessionsID...).Err(); err != nil {
		return err
	}

	for _, sessionID := range sessionsID {
		delete(AccountSessions, sessionID)
	}

	if err := r.UpdateSessionsForAccount(ctx, AccountSessionsCache{AccountSessions}, accountID); err != nil {
		return err
	}

	return nil
}

func (r *redisSessionsCache) cacheAccountSession(ctx context.Context, toCache model.SessionCache) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.cacheAccountSession")
	defer span.Finish()

	body, err := r.account_sessions_rdb.Get(ctx, toCache.AccountID).Bytes()
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

	return r.UpdateSessionsForAccount(ctx, sessionsCache, toCache.AccountID)
}

func (r *redisSessionsCache) GetSessionCache(ctx context.Context, sessionID string) (model.SessionCache, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.GetSessionCache")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	body, err := r.sessions_rdb.Get(ctx, sessionID).Bytes()
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

func (r *redisSessionsCache) UpdateLastActivityForSession(ctx context.Context,
	cachedSession model.SessionCache, sessionID string, LastActivityTime time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"SessionsCache.UpdateLastActivityForSession")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	cachedSession.LastActivity = LastActivityTime
	return r.CacheSession(ctx, cachedSession)
}

type sessionInfo struct {
	ClientIP     string    `json:"client_ip"`
	SessionInfo  string    `json:"session_info"` // like device or browser
	LastActivity time.Time `json:"last_activity"`
}

type AccountSessionsCache struct {
	Sessions map[string]sessionInfo `json:"sessions"`
}

func (r *redisSessionsCache) GetSessionsForAccount(ctx context.Context, accountID string) (map[string]sessionInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.GetSessionsForAccount")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	body, err := r.account_sessions_rdb.Get(ctx, accountID).Bytes()
	if err == redis.Nil {
		return map[string]sessionInfo{}, nil // not error, just empty sessions map for account
	}

	r.logger.Debug("Unmarshal cache data")
	var sessions AccountSessionsCache
	err = json.Unmarshal(body, &sessions)

	return sessions.Sessions, err
}

func (r *redisSessionsCache) UpdateSessionsForAccount(ctx context.Context,
	sessions AccountSessionsCache, accountID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.UpdateSessionsForAccount")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	if len(sessions.Sessions) == 0 {
		err := r.account_sessions_rdb.Del(ctx, accountID).Err()
		if err == redis.Nil {
			return nil
		}
		return err
	}

	r.logger.Debug("Marshalling data")
	serialized, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	r.logger.Debug("Caching data")
	_, err = r.account_sessions_rdb.Set(ctx, accountID, serialized, r.SessionTTL).Result()

	return err
}
