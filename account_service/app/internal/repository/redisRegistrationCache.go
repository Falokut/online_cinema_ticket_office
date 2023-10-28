package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type redisRegistrationCache struct {
	rdb    *redis.Client
	logger logging.Logger
}

func NewRedisRegistrationCache(opt *redis.Options, logger logging.Logger) (*redisRegistrationCache, error) {
	logger.Info("Creating registration cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Info("Pinging registration cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	return &redisRegistrationCache{rdb: rdb, logger: logger}, nil
}

func (r *redisRegistrationCache) ShutDown() {
	r.logger.Info("Registration cache repository shutting down")
	r.rdb.Close()
}

func (r redisRegistrationCache) IsAccountInCache(ctx context.Context, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsCache.UpdateSessionsForAccount")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	num, err := r.rdb.Exists(ctx, email).Result()
	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func (r *redisRegistrationCache) CacheAccount(ctx context.Context,
	email string, account CachedAccount, NonActivatedAccountTTL time.Duration) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.CacheAccount")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	r.logger.Info("Marshalling data")
	serialized, err := json.Marshal(&account)
	if err != nil {
		return err
	}

	_, err = r.rdb.Set(ctx, email, serialized, NonActivatedAccountTTL).Result()
	if err != nil {
		return errors.New("Can't cache account")
	}

	return nil
}

func (r redisRegistrationCache) GetCachedAccount(ctx context.Context, email string) (CachedAccount, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.GetCachedAccount")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	r.logger.Printf("Getting account with %s email", email)
	body, err := r.rdb.Get(ctx, email).Bytes()
	if err != nil && err != redis.Nil {
		return CachedAccount{}, err
	} else if err == redis.Nil {
		return CachedAccount{}, errors.New("Account not found")
	}

	var account CachedAccount
	err = json.Unmarshal(body, &account)

	return account, err
}

func (r *redisRegistrationCache) DeleteAccountFromCache(ctx context.Context, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationCache.DeleteAccountFromCache")
	span.SetTag("custom-tag", "redis cache")
	defer span.Finish()

	return r.rdb.Del(ctx, email).Err()
}
