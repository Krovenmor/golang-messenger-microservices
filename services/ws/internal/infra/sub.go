package infra

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis"
	"context"
	"fmt"
	"log"
)

type RedisSubscriber struct {
	sub               *redis.RedisSubscriber
	userPattern       string
	chatPattern       string
	userStatusPattern string

	chatReader *chatsReader
}

func NewSubscriber(sub *redis.RedisSubscriber, conf *config.RedisChannelsConfig) *RedisSubscriber {
	rSub := &RedisSubscriber{
		sub:               sub,
		userPattern:       conf.UserEventsPattern,
		chatPattern:       conf.ChatEventsPattern,
		userStatusPattern: conf.UserStatusPattern,
	}
	rSub.chatReader = newChatsReader(rSub)
	return rSub
}

func (s *RedisSubscriber) SubscribeOnUserEvents(ctx context.Context, userId string) (<-chan []byte, func(), error) {
	channel := fmt.Sprintf(s.userPattern, userId)
	log.Printf("Subscribed to a user channel: %q, id: %q", channel, userId)
	return s.sub.Subscribe(ctx, channel)
}

func (s *RedisSubscriber) SubscribeOnChatEventsInternal(ctx context.Context, chatId string) (<-chan []byte, func(), error) {
	channel := fmt.Sprintf(s.chatPattern, chatId)
	log.Printf("Subscribed to a chat channel: %q, id: %q", channel, chatId)
	return s.sub.Subscribe(ctx, channel)
}

func (s *RedisSubscriber) SubscribeOnChatEvents(ctx context.Context, chatId string) (<-chan []byte, func(), error) {
	log.Printf("Trying to subscribing to a chat, id: %q", chatId)
	return s.chatReader.Subscribe(chatId)
}

func (s *RedisSubscriber) SubscribeOnUserStatuses(ctx context.Context, userId string) (<-chan []byte, func(), error) {
	channel := fmt.Sprintf(s.userStatusPattern, userId)
	log.Printf("Subscribed to a user status channel: %q, id: %q", channel, userId)
	return s.sub.Subscribe(ctx, channel)
}
