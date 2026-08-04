package infra

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/redis"
	"MyMessenger/services/ws/internal/config"
	"context"
	"fmt"
	"log"
)

type RedisPublisher struct {
	pub     *redis.RedisPublisher
	pattern string
}

func NewPublisher(pub *redis.RedisPublisher, conf *config.RedisPatternConfig) *RedisPublisher {
	return &RedisPublisher{
		pub:     pub,
		pattern: conf.PatternPub,
	}
}

func (p *RedisPublisher) PublishUserStatus(ctx context.Context, status broker.StatusEvent) {
	channel := fmt.Sprintf(p.pattern, status.UserId)
	err := p.pub.PublishEvent(ctx, fmt.Sprintf(p.pattern, status.UserId), status)
	log.Printf("Published new status, channel: %q, event: %v, err: %v", channel, status, err)
}
