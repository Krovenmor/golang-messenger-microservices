package infra

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"
	"fmt"
	"log"
)

type RedisPublisher struct {
	pub               *redis.RedisPublisher
	userStatusPattern string
	banChannel        string
}

func NewPublisher(pub *redis.RedisPublisher, conf *config.RedisChannelsConfig) *RedisPublisher {
	return &RedisPublisher{
		pub:               pub,
		userStatusPattern: conf.UserStatusPattern,
		banChannel:        conf.UserBanChannel,
	}
}

func (p *RedisPublisher) PublishUserStatus(ctx context.Context, status broker.StatusPayload) {
	channel := fmt.Sprintf(p.userStatusPattern, status.UserId)
	event := broker.Event{
		Type:    broker.StatusEvent,
		Payload: status,
	}
	err := p.pub.PublishEvent(ctx, channel, event)
	log.Printf("Published new status, channel: %q, event: %v, err: %v", channel, event, err)
}

func (p *RedisPublisher) PublishBanEvent(ctx context.Context, payload broker.BanEventPayload) {
	event := broker.Event{
		Type:    broker.BanEvent,
		Payload: payload,
	}
	err := p.pub.PublishEvent(ctx, p.banChannel, event)
	log.Printf("Published new ban, channel: %q, event: %v, err: %v", p.banChannel, event, err)
}
