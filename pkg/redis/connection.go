package redis

import (
	"MyMessenger/pkg/config"
	"context"
	"log"
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
		log.Printf("NewRedisSubscriber: Trouble with connecting to redis server, err: %q", err.Error())
		return nil, err
	}

	return rdClient, nil
}
