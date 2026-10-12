package redis

import (
	"MyMessenger/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var waitTime = time.Second * 3

func NewRedisClient(conf *config.RedisConfig) (*redis.Client, error) {
	rdClient := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	err := rdClient.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return rdClient, nil
}
