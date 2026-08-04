package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	rdClient *redis.Client
}

func NewRedisSubscriber(rdClient *redis.Client) *RedisSubscriber {
	return &RedisSubscriber{
		rdClient: rdClient,
	}
}

func (s *RedisSubscriber) Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error) {
	sub := s.rdClient.Subscribe(ctx, channel)

	_, err := sub.Receive(ctx)
	if err != nil {
		log.Printf("Subscribe: Failed to sub.Receive, err: %q", err.Error())
		return nil, nil, err
	}

	out := make(chan []byte, 10)

	go func() {
		defer close(out)
		ch := sub.Channel()
		for msg := range ch {
			select {
			case out <- []byte(msg.Payload):
			case <-ctx.Done():
				return
			}
		}
	}()

	cancel := func() {
		err := sub.Close()
		if err != nil {
			log.Printf("Subscribe: Failed to close pub, err: %q", err.Error())
		}
	}

	return out, cancel, nil
}
