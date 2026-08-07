package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/redis"
	"MyMessenger/services/msg/internal/config"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type publisher struct {
	pub            *redis.RedisPublisher
	pubUserPattern string
	pubChatPattern string
}

func NewPublisher(pub *redis.RedisPublisher, conf *config.RedisPatternConfig) (*publisher, error) {
	return &publisher{
		pub:            pub,
		pubUserPattern: conf.PubUserPattern,
		pubChatPattern: conf.PubChatPattern,
	}, nil
}

func (p *publisher) toChat(chatId uuid.UUID) string {
	return fmt.Sprintf(p.pubChatPattern, chatId.String())
}

func (p *publisher) PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewChatType,
		Payload: broker.NewChatPayload{
			ChatId: chatId.String(),
		},
	}
	redis.PublishEventToGroup(ctx, p.pub, p.pubUserPattern, event, usersTo)
}

func (p *publisher) PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.NewMessageType,
		Payload: broker.NewMessagePayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.pub.PublishEvent(ctx, p.toChat(chatId), event)
}

func (p *publisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.pub.PublishEvent(ctx, p.toChat(chatId), event)
}

func (p *publisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.pub.PublishEvent(ctx, p.toChat(chatId), event)
}
