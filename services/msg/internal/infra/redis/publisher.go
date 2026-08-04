package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/redis"
	"context"

	"github.com/google/uuid"
)

var pattern = "user:%s:events"

type publisher struct {
	pub *redis.RedisPublisher
}

func NewPublisher(pub *redis.RedisPublisher) (*publisher, error) {
	return &publisher{
		pub: pub,
	}, nil
}

func (p *publisher) PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewChatType,
		Payload: broker.NewChatPayload{
			ChatId: chatId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, pattern, event, usersTo)
}

func (p *publisher) PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewMessageType,
		Payload: broker.NewMessagePayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, pattern, event, usersTo)
}

func (p *publisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, pattern, event, usersTo)
}

func (p *publisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, pattern, event, usersTo)
}
