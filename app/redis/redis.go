package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

type Service interface {
	GetClient(id string) (string, error)
}

type Impl struct {
	redisClient *redis.Client
}

func NewRedisService() (*Impl, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Impl{redisClient: redisClient}, nil
}

func (r Impl) GetClient(id string) (string, error) {
	prefix := os.Getenv("REDIS_PREFIX")
	key := prefix + "users" + id

	return r.redisClient.Get(context.Background(), key).Result()
}
