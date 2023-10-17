package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/account_service/pkg/logging"
	"github.com/redis/go-redis/v9"
)

type redisRegistrationCache struct {
	rdb    *redis.Client
	logger logging.Logger
}

func NewRedisRegistrationCache(opt *redis.Options, logger logging.Logger) (*redisRegistrationCache, error) {
	logger.Println("Creating registration cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("Can't create new redis client")
	}

	logger.Println("Pinging registration cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection is not established: %s", err.Error()))
	}

	return &redisRegistrationCache{rdb: rdb, logger: logger}, nil
}

func (r *redisRegistrationCache) ShutDown() {
	r.logger.Println("Registration cache repository shutting down")
	r.rdb.Close()
}

func (r redisRegistrationCache) IsAccountInCache(email string) (bool, error) {
	num, err := r.rdb.Exists(context.Background(), email).Result()
	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func (r *redisRegistrationCache) CacheAccount(email string, account CachedAccount, NonActivatedAccountTTL time.Duration) error {
	r.logger.Println("Marshalling data")
	serialized, err := json.Marshal(&account)
	if err != nil {
		return err
	}

	_, err = r.rdb.Set(context.Background(), email, serialized, NonActivatedAccountTTL).Result()
	if err != nil {
		return errors.New("Can't cache account")
	}

	return nil
}

func (r redisRegistrationCache) GetCachedAccount(email string) (CachedAccount, error) {
	r.logger.Printf("Getting account with %s email", email)
	body, err := r.rdb.Get(context.Background(), email).Bytes()
	if err != nil && err != redis.Nil {
		return CachedAccount{}, err
	} else if err == redis.Nil {
		return CachedAccount{}, errors.New("Account not found")
	}

	var account CachedAccount
	err = json.Unmarshal(body, &account)
	if err != nil {
		return CachedAccount{}, err
	}

	return account, nil
}

func (r *redisRegistrationCache) DeleteAccountFromCache(email string) error {
	return r.rdb.Del(context.Background(), email).Err()
}
