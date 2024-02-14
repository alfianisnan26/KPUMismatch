package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kawalrealcount/internal/data/dao"
	"time"
)

type repo struct {
	client *redis.Client
}

type Param struct {
	Host string
}

func New(param Param) (dao.Cache, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr: param.Host,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &repo{
		client: rdb,
	}, nil
}

func (repo repo) Get(key string, expiry time.Duration, fallback func(string) ([]byte, error)) ([]byte, error) {
	ctx := context.Background()

	// Try to get the value from Redis cache
	val, err := repo.client.Get(ctx, key).Bytes()
	if err == nil {
		// Value found in cache, return it
		return val, nil
	}

	if err == redis.Nil {
		// Key not found in cache, fetch from the fallback function
		fallbackVal, err := fallback(key)
		if err != nil {
			return nil, err
		}

		// Store the value in cache with a TTL of 1 hour
		err = repo.client.Set(ctx, key, fallbackVal, expiry).Err()
		if err != nil {
			return nil, err
		}

		return fallbackVal, nil
	}

	// Some other error occurred
	return nil, err
}
