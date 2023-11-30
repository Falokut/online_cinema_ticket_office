package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type redisRegistrationCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
}

func (r *redisRegistrationCache) PingContext(ctx context.Context) error {
	if err := r.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging registration cache: %w", err)
	}

	return nil
}

// NewRedisRegistrationCache initializes a new instance of redisRegistrationCache with the provided options and logger.
func NewRedisRegistrationCache(opt *redis.Options, logger *logrus.Logger) (*redisRegistrationCache, error) {
	logger.Info("Creating registration cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging registration cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &redisRegistrationCache{rdb: rdb, logger: logger}, nil
}

// Shutdown gracefully shuts down the registration cache repository.
func (r *redisRegistrationCache) Shutdown() error {
	r.logger.Info("Registration cache repository shutting down")
	return r.rdb.Close()
}

// IsAccountInCache checks if the provided email account is present in the cache.
func (r redisRegistrationCache) IsAccountInCache(ctx context.Context, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.UpdateSessionsForAccount")
	defer span.Finish()

	num, err := r.rdb.Exists(ctx, email).Result()
	if err != nil {
		return false, err
	}

	return num > 0, nil
}

// CacheAccount caches the provided account information with the specified email in the Redis cache.
// It marshals the account data into JSON and sets it in the cache with the specified TTL.
func (r *redisRegistrationCache) CacheAccount(ctx context.Context,
	email string, account CachedAccount, NonActivatedAccountTTL time.Duration) error {
	// Start a new span for caching the account information
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.CacheAccount")
	defer span.Finish()

	// Log the marshalling of data
	r.logger.Info("Marshalling data")
	serialized, err := json.Marshal(&account)
	if err != nil {
		return err
	}

	// Set the serialized account data in the Redis cache with the specified TTL
	_, err = r.rdb.Set(ctx, email, serialized, NonActivatedAccountTTL).Result()
	if err != nil {
		return fmt.Errorf("can't cache account %s", err.Error())
	}

	return nil
}

// GetCachedAccount retrieves the cached account information for the specified email from the Redis cache.
// It returns the cached account data and any encountered error during retrieval.
func (r redisRegistrationCache) GetCachedAccount(ctx context.Context, email string) (CachedAccount, error) {
	// Start a new span for retrieving the cached account information
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.GetCachedAccount")
	defer span.Finish()

	// Log the retrieval of account with the specified email
	r.logger.Debugf("Getting account with %s email", email)
	body, err := r.rdb.Get(ctx, email).Bytes()
	if err != nil {
		return CachedAccount{}, err
	}

	// Unmarshal the retrieved serialized data into a CachedAccount struct
	var account CachedAccount
	err = json.Unmarshal(body, &account)

	return account, err
}

// DeleteAccountFromCache deletes the account information associated with the specified email from the Redis cache.
func (r *redisRegistrationCache) DeleteAccountFromCache(ctx context.Context, email string) error {
	// Start a new span for deleting the account information from the cache
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.DeleteAccountFromCache")
	defer span.Finish()

	// Delete the account information from the Redis cache
	return r.rdb.Del(ctx, email).Err()
}
