package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"

	"github.com/google/uuid"
)

type redisStreamsPublisher struct {
	pub    *redis.RedisStreamsPublisher
	stream string
}

func NewRedisStreamsPublisher(pub *redis.RedisStreamsPublisher, conf *config.RedisChannelsConfig) *redisStreamsPublisher {
	return &redisStreamsPublisher{pub: pub, stream: conf.ProfileStream}
}

func (r *redisStreamsPublisher) PubNewProfile(ctx context.Context, userId uuid.UUID) error {
	return r.pub.PublishEvent(ctx, r.stream, broker.ProfileEventDTO{
		Type: broker.NewProfileEvent,
		Payload: broker.ProfilePayload{
			UserId: userId.String(),
		},
	})
}
