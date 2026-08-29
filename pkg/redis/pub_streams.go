package redis

import (
	"MyMessenger/pkg/config"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisStreamsPublisher struct {
	rdClient *redis.Client
	maxLen   int64
}

func NewRedisStreamsPublisher(rdClient *redis.Client, conf *config.RedisConfig) *RedisStreamsPublisher {
	return &RedisStreamsPublisher{
		rdClient: rdClient,
		maxLen:   conf.StreamsMaxLen,
	}
}

func (p *RedisStreamsPublisher) PublishEvent(ctx context.Context, stream string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.rdClient.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: p.maxLen,
		Approx: true,
		Values: map[string]any{
			"event": data,
		}}).Err()
}
