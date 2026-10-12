package redis

import (
	"MyMessenger/pkg/logger"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisSubscriber struct {
	rdClient *redis.Client
	log      *logger.Logger
}

func NewRedisSubscriber(rdClient *redis.Client, log *logger.Logger) *RedisSubscriber {
	return &RedisSubscriber{
		rdClient: rdClient,
		log:      log,
	}
}

func (s *RedisSubscriber) subscribe(ctx context.Context, sub *redis.PubSub) (<-chan []byte, func(), error) {
	_, err := sub.Receive(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("subscribe receive: %w", err)
	}

	out := make(chan []byte, 10)

	go func() {
		defer close(out)
		ch := sub.Channel()
		for msg := range ch {
			select {
			case out <- []byte(msg.Payload):
			case <-ctx.Done():
				return
			}
		}
	}()

	cancel := func() {
		err := sub.Close()
		if err != nil {
			s.log.Warn("failed to close subscription", zap.Error(err))
		}
	}

	return out, cancel, nil
}

func (s *RedisSubscriber) PatternSubscribe(ctx context.Context, pattern string) (<-chan []byte, func(), error) {
	sub := s.rdClient.PSubscribe(ctx, pattern)
	return s.subscribe(ctx, sub)
}

func (s *RedisSubscriber) Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error) {
	sub := s.rdClient.Subscribe(ctx, channel)
	return s.subscribe(ctx, sub)
}
