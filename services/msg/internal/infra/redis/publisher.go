package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type publisher struct {
	pub            *redis.RedisPublisher
	pubUserPattern string
	pubChatPattern string
}

func NewPublisher(pub *redis.RedisPublisher, conf *config.RedisChannelsConfig) (*publisher, error) {
	return &publisher{
		pub:            pub,
		pubUserPattern: conf.UserEventsPattern,
		pubChatPattern: conf.ChatEventsPattern,
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
	err := redis.PublishEventToGroup(ctx, p.pub, p.pubUserPattern, event, usersTo)
	if err != nil {
		log.Printf("PublishNewChat: trouble with PublishEventToGroup, err: %q", err)
	}
}

func (p *publisher) PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.NewMessageType,
		Payload: broker.NewMessagePayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	err := p.pub.PublishEvent(ctx, p.toChat(chatId), event)
	if err != nil {
		log.Printf("PublishNewMessage: trouble with PublishEvent, err: %q", err)
	}
}

func (p *publisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	err := p.pub.PublishEvent(ctx, p.toChat(chatId), event)
	if err != nil {
		log.Printf("PublishMessageWasRedacted: trouble with PublishEvent, err: %q", err)
	}
}

func (p *publisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	err := p.pub.PublishEvent(ctx, p.toChat(chatId), event)
	if err != nil {
		log.Printf("PublishMessageWasRedacted: trouble with PublishEvent, err: %q", err)
	}
}
