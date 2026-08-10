package pub

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"
)

type redisPublisher struct {
	pub        *redis.RedisPublisher
	pubChannel string
}

func NewRedisPublisher(pub *redis.RedisPublisher, conf *config.RedisChannelsConfig) *redisPublisher {
	return &redisPublisher{pub: pub, pubChannel: conf.UserBanEventChannel}
}

func (p *redisPublisher) Pub(ctx context.Context, payload broker.BanEventPayload) error {
	return p.pub.PublishEvent(ctx, p.pubChannel, broker.Event{
		Type:    broker.BanEvent,
		Payload: payload,
	})
}
