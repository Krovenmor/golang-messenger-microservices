package redis

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	rdClient *redis.Client
}

func NewRedisPublisher(conf *config.RedisConfig) (*RedisPublisher, error) {
	rdClient := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := rdClient.Ping(ctx).Err()
	if err != nil {
		log.Printf("Trouble with connecting to redis server, err: %q", err.Error())
		return nil, err
	}

	return &RedisPublisher{
		rdClient: rdClient,
	}, nil
}

func (p *RedisPublisher) publishEvent(ctx context.Context, event broker.Event, users []uuid.UUID) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("publishEvent: trouble with Marshaling event:%v, err:%q", event, err.Error())
		return
	}
	for _, user := range users {
		channel := fmt.Sprintf("user:%s:events", user.String())
		err := p.rdClient.Publish(ctx, channel, data).Err()
		if err != nil {
			log.Printf("publishEvent: trouble with Publishing event:%v, err:%q, to:%q", event, err.Error(), user.String())
		}
	}
}

func (p *RedisPublisher) PublishNewChat(ctx context.Context, chatId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewChatType,
		Payload: broker.NewChatPayload{
			ChatId: chatId.String(),
		},
	}
	p.publishEvent(ctx, event, usersTo)
}

func (p *RedisPublisher) PublishNewMessage(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.NewMessageType,
		Payload: broker.NewMessagePayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.publishEvent(ctx, event, usersTo)
}

func (p *RedisPublisher) PublishMessageWasRedacted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageRedactedType,
		Payload: broker.MessageRedactedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.publishEvent(ctx, event, usersTo)
}

func (p *RedisPublisher) PublishMessageWasDeleted(ctx context.Context, chatId, msgId uuid.UUID, usersTo []uuid.UUID) {
	event := broker.Event{
		Type: broker.MessageDeletedType,
		Payload: broker.MessageDeletedPayload{
			ChatId: chatId.String(),
			MsgId:  msgId.String(),
		},
	}
	p.publishEvent(ctx, event, usersTo)
}
