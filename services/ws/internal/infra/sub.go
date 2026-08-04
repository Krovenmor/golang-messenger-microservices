package infra

import (
	"MyMessenger/pkg/redis"
	"MyMessenger/services/ws/internal/config"
	"context"
	"fmt"
	"log"
)

type RedisSubscriber struct {
	sub     *redis.RedisSubscriber
	pattern string
}

func NewSubscriber(sub *redis.RedisSubscriber, conf *config.RedisPatternConfig) *RedisSubscriber {
	return &RedisSubscriber{
		sub:     sub,
		pattern: conf.PatternSub,
	}
}

func (s *RedisSubscriber) Subscribe(ctx context.Context, userId string) (<-chan []byte, func(), error) {
	channel := fmt.Sprintf(s.pattern, userId)
	log.Printf("Subscribed to a channel: %q, id: %q", channel, userId)
	return s.sub.Subscribe(ctx, channel)
}
