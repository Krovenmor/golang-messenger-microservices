package pub

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"
)

type publisher struct {
	pub     *redis.RedisPublisher
	channel string
}

func NewPublisher(pub *redis.RedisPublisher, channels *config.RedisChannelsConfig) *publisher {
	return &publisher{pub: pub, channel: channels.UserVerificationChannel}
}

func (p *publisher) PublishEmailVerification(ctx context.Context, email, code string) error {
	event := broker.Event{
		Type: broker.EmailVerificationType,
		Payload: broker.EmailVerificationPayload{
			Email: email,
			Code:  code,
		},
	}
	return p.pub.PublishEvent(ctx, p.channel, event)
}
