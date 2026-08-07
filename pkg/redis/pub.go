package redis

import (
	"MyMessenger/pkg/broker"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	rdClient *redis.Client
}

func NewRedisPublisher(rdClient *redis.Client) *RedisPublisher {
	return &RedisPublisher{
		rdClient: rdClient,
	}
}

func (p *RedisPublisher) PublishEvent(ctx context.Context, channel string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("PublishEvent: trouble with Marshaling event:%v, err:%q", event, err.Error())
		return err
	}
	err = p.rdClient.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Printf("PublishEvent: trouble with Publishing event:%v, err:%q, to:%q", event, err.Error(), channel)
		return err
	}
	log.Printf("PublishEvent: Publish to ch: %q", channel)
	return nil
}

type WithString interface {
	String() string
}

func PublishEventToGroup[T WithString](ctx context.Context, p *RedisPublisher, pattern string, event broker.Event, group []T) error {
	if len(group) == 0 {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("PublishEventToGroup: trouble with Marshaling event: %v, err: %v", event, err)
		return err
	}

	_, err = p.rdClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, arg := range group {
			channel := fmt.Sprintf(pattern, arg.String())
			err := pipe.Publish(ctx, channel, data).Err()
			if err != nil {
				log.Printf("PublishEventToGroup: Publish err:%q, to ch: %q", err, channel)
			} else {
				log.Printf("PublishEventToGroup: Publish to ch: %q", channel)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("PublishEventToGroup: pipeline publish error for pattern %q: %v", pattern, err)
		return err
	}

	return nil
}
