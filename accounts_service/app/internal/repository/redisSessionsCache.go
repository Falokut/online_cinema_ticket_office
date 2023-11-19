package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type redisSessionsCache struct {
	sessions_rdb         *redis.Client
	account_sessions_rdb *redis.Client
	logger               logging.Logger
	SessionTTL           time.Duration
}

// NewSessionCache creates a new session cache using the provided Redis options, logger, and session TTL.
// It initializes two Redis clients for session and account session caching, and verifies the connection to each Redis instance.
func NewSessionCache(sessionCacheOpt *redis.Options, accountSessionsOpt *redis.Options, logger logging.Logger, SessionTTL time.Duration) (*redisSessionsCache, error) {
	logger.Infoln("Creating session cache client")

	// Initialize a Redis client for session caching
	sessions_rdb := redis.NewClient(sessionCacheOpt)
	if sessions_rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Infoln("Pinging session cache client")
	// Verify the connection to the session cache Redis instance
	_, err := sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	logger.Infoln("Creating account sessions cache client")
	// Initialize a Redis client for account sessions caching
	account_sessions_rdb := redis.NewClient(accountSessionsOpt)
	if account_sessions_rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Infoln("Pinging session cache client")
	// Verify the connection to the account sessions cache Redis instance
	_, err = account_sessions_rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	// Return the initialized redisSessionsCache with the configured Redis clients, logger, and session TTL
	return &redisSessionsCache{sessions_rdb: sessions_rdb, account_sessions_rdb: account_sessions_rdb, logger: logger, SessionTTL: SessionTTL}, nil
}

// Shutdown gracefully shuts down the token cache repository by closing the Redis client for session caching.
func (r *redisSessionsCache) Shutdown() error {
	r.logger.Infoln("Token cache repository shutting down")
	return r.sessions_rdb.Close()
}

func (r *redisSessionsCache) CacheSession(ctx context.Context, toCache model.SessionCache) error {
	// Create a new span for caching the session data
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.CacheSession")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// Log a message indicating the marshalling of the data
	r.logger.Info("Marshalling data")
	// Marshal the session data into a JSON format
	serialized, err := json.Marshal(toCache)
	if err != nil {
		return err
	}

	// Log a message indicating the caching of the session data
	r.logger.Info("Caching sessions data")
	// Cache the serialized session data in Redis with the specified TTL
	_, err = r.sessions_rdb.Set(ctx, toCache.SessionID, serialized, r.SessionTTL).Result()
	if err != nil {
		return err
	}

	// Cache the account session data
	if err := r.cacheAccountSession(ctx, toCache); err != nil {
		// If an error occurs, delete the session data from the cache
		r.sessions_rdb.Del(ctx, toCache.SessionID)
		return err
	}

	return nil
}

func (r *redisSessionsCache) TerminateSessions(ctx context.Context, sessionsID []string, accountID string) error {
	// Create a new span for terminating sessions
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.TerminateSessions")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// Get the list of sessions associated with the specified account
	AccountSessions, err := r.GetSessionsForAccount(ctx, accountID)
	if err != nil {
		return err
	}

	// Iterate through the session IDs and delete them from the cache
	if err := r.sessions_rdb.Del(ctx, sessionsID...).Err(); err != nil {
		// Handle the error
		return err
	}

	// Update the account's remaining sessions
	if err := r.UpdateSessionsForAccount(ctx, AccountSessionsCache{AccountSessions}, accountID); err != nil {
		return err
	}

	return nil
}

// cacheAccountSession caches the session information for a specific account in the Redis cache.
// It retrieves the existing session data, updates it with new information, and then updates the cache.
// This function is a part of the RedisSessionsCache struct.
func (r *redisSessionsCache) cacheAccountSession(ctx context.Context, toCache model.SessionCache) error {
	// Start a new span for tracing the cacheAccountSession operation
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.cacheAccountSession")
	defer span.Finish()

	// Attempt to retrieve the existing session data for the account from the Redis cache
	body, err := r.account_sessions_rdb.Get(ctx, toCache.AccountID).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}

	// Initialize a struct to hold the session cache data
	var sessionsCache AccountSessionsCache
	sessionsCache.Sessions = make(map[string]sessionInfo)

	// Unmarshal the retrieved data if it exists
	if err != redis.Nil {
		// Log an informational message indicating unmarshaling of data
		r.logger.Info("Unmarshal data")
		if err := json.Unmarshal(body, &sessionsCache); err != nil {
			return err
		}
	}

	// Update the session information with the new data
	sessionsCache.Sessions[toCache.SessionID] = sessionInfo{ClientIP: toCache.ClientIP,
		SessionInfo: toCache.SessionInfo, LastActivity: toCache.LastActivity}

	// Update the sessions for the account in the Redis cache
	return r.UpdateSessionsForAccount(ctx, sessionsCache, toCache.AccountID)
}

// GetSessionCache retrieves the session cache for a given session ID from the Redis cache.
// It starts a new span for tracing, retrieves the session data from the cache, and unmarshals it into a model.SessionCache.
// If the session data is not found in the cache, it returns an error indicating that the session was not found.
func (r *redisSessionsCache) GetSessionCache(ctx context.Context, sessionID string) (model.SessionCache, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.GetSessionCache")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// Retrieve the session data for the given session ID from the Redis cache
	body, err := r.sessions_rdb.Get(ctx, sessionID).Bytes()
	if err == redis.Nil {
		return model.SessionCache{}, ErrSessionNotFound
	}

	// Unmarshal the retrieved cache data into a model.SessionCache struct
	r.logger.Info("Unmarshal cache data")
	var session model.SessionCache
	if err := json.Unmarshal(body, &session); err != nil {
		return model.SessionCache{}, err
	}

	return session, nil
}

// UpdateLastActivityForSession updates the last activity time for a cached session in the Redis cache.
// It starts a new span for tracing, updates the LastActivity field of the cached session, and then caches the updated session.
func (r *redisSessionsCache) UpdateLastActivityForSession(ctx context.Context,
	cachedSession model.SessionCache, sessionID string, LastActivityTime time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.UpdateLastActivityForSession")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// Update the LastActivity field of the cached session with the provided LastActivityTime
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

// GetSessionsForAccount retrieves the sessions associated with the specified account from the Redis cache.
func (r *redisSessionsCache) GetSessionsForAccount(ctx context.Context, accountID string) (map[string]sessionInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.GetSessionsForAccount")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// Retrieve the cached data for the specified accountID from the Redis database.
	body, err := r.account_sessions_rdb.Get(ctx, accountID).Bytes()
	if err == redis.Nil {
		return map[string]sessionInfo{}, nil // No error, just an empty sessions map for the account
	}

	// Unmarshal the cached data into the AccountSessionsCache struct.
	r.logger.Info("Unmarshal cache data")
	var sessions AccountSessionsCache
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		return nil, err
	}

	return sessions.Sessions, nil
}

// UpdateSessionsForAccount updates the sessions associated with the specified account in the Redis cache.
func (r *redisSessionsCache) UpdateSessionsForAccount(ctx context.Context, sessions AccountSessionsCache, accountID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.UpdateSessionsForAccount")
	span.SetTag("database", "redis cache")
	defer span.Finish()

	// If the sessions map is empty, remove the corresponding entry from the Redis database.
	if len(sessions.Sessions) == 0 {
		err := r.account_sessions_rdb.Del(ctx, accountID).Err()
		if err == redis.Nil {
			return nil
		}
		return err
	}

	// Marshal the sessions data into a JSON object.
	r.logger.Info("Marshalling data")
	serialized, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	// Store the serialized sessions data in the Redis database with the specified accountID.
	r.logger.Info("Caching data")
	_, err = r.account_sessions_rdb.Set(ctx, accountID, serialized, r.SessionTTL).Result()

	return err
}
