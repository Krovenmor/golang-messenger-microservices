package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/redis"
	"MyMessenger/services/msg/internal/config"
	"context"

	"github.com/google/uuid"
)

type publisher struct {
	pub        *redis.RedisPublisher
	pubPattern string
}

func NewPublisher(pub *redis.RedisPublisher, conf *config.RedisPatternConfig) (*publisher, error) {
	return &publisher{
		pub:        pub,
		pubPattern: conf.PubPattern,
	}, nil
}

func (p *publisher) PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewChatType,
		Payload: broker.NewChatPayload{
			ChatId: chatId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, p.pubPattern, event, usersTo)
}

func (p *publisher) PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewMessageType,
		Payload: broker.NewMessagePayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, p.pubPattern, event, usersTo)
}

func (p *publisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, p.pubPattern, event, usersTo)
}

func (p *publisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, p.pubPattern, event, usersTo)
}
