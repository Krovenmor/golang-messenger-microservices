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

func (p *publisher) publishToChat(ctx context.Context, chatId uuid.UUID, event broker.Event) {
	err := p.pub.PublishEvent(ctx, p.toChat(chatId), event)
	if err != nil {
		log.Printf("publishToChat: trouble with PublishEvent to chat:%v, event:%v, err: %v", chatId, event, err)
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
	p.publishToChat(ctx, chatId, event)
}

func (p *publisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.publishToChat(ctx, chatId, event)
}

func (p *publisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.publishToChat(ctx, chatId, event)
}

func (p *publisher) PublishNewReaction(ctx context.Context, chatId, msgId, userId uuid.UUID, emoji string) {
	event := broker.Event{
		Type: broker.NewReactionType,
		Payload: broker.NewReactionPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
			UserId: userId.String(),
			Emoji:  emoji,
		},
	}
	p.publishToChat(ctx, chatId, event)
}

func (p *publisher) PublishDelReaction(ctx context.Context, chatId, msgId, userId uuid.UUID, emoji string) {
	event := broker.Event{
		Type: broker.DelReactionType,
		Payload: broker.DelReactionPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
			UserId: userId.String(),
			Emoji:  emoji,
		},
	}
	p.publishToChat(ctx, chatId, event)
}
