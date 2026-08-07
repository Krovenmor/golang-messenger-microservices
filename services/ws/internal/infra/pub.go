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
		pattern: conf.PubPattern,
	}
}

func (p *RedisPublisher) PublishUserStatus(ctx context.Context, status broker.StatusPayload) {
	channel := fmt.Sprintf(p.pattern, status.UserId)
	event := broker.Event{
		Type:    broker.StatusEvent,
		Payload: status,
	}
	err := p.pub.PublishEvent(ctx, channel, event)
	log.Printf("Published new status, channel: %q, event: %v, err: %v", channel, event, err)
}
