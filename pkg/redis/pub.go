package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisPublisher struct {
	rdClient *redis.Client
	log      *logger.Logger
}

func NewRedisPublisher(rdClient *redis.Client, log *logger.Logger) *RedisPublisher {
	return &RedisPublisher{
		rdClient: rdClient,
		log:      log,
	}
}

func (p *RedisPublisher) PublishEvent(ctx context.Context, channel string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	err = p.rdClient.Publish(ctx, channel, data).Err()
	if err != nil {
		return fmt.Errorf("publish to %q: %w", channel, err)
	}
	p.log.Ctx(ctx).Debug("event published", zap.String("channel", channel))
	return nil
}

type WithString interface {
	String() string
}

func PublishEventToGroup[T WithString](ctx context.Context, p *RedisPublisher, pattern string, event broker.Event, group []T) error {
	if len(group) == 0 {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	_, err = p.rdClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, arg := range group {
			pipe.Publish(ctx, fmt.Sprintf(pattern, arg.String()), data)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("publish to group %q: %w", pattern, err)
	}

	p.log.Ctx(ctx).Debug("event published to group", zap.String("pattern", pattern), zap.Int("size", len(group)))
	return nil
}
