package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CallbackFunc = func(ctx context.Context, msg redis.XMessage)

type RedisStreamsConsumer struct {
	rdClient *redis.Client

	stream   string
	group    string
	consumer string
}

func GenUniqueConsumerName(serviceName string) string {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	return fmt.Sprintf("%s-%s-%d", serviceName, host, os.Getpid())
}

func NewRedisStreamsConsumer(rdClient *redis.Client, stream, group, consumer string) *RedisStreamsConsumer {
	return &RedisStreamsConsumer{
		rdClient: rdClient,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (r *RedisStreamsConsumer) InitGroup(ctx context.Context) error {
	err := r.rdClient.XGroupCreateMkStream(ctx, r.stream, r.group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Printf("RedisStreamsConsumer: trouble with XGroupCreateMkStream, err: %q", err)
		return err
	}
	return nil
}

func (r *RedisStreamsConsumer) Start(ctx context.Context, callback CallbackFunc) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer [%s] stopped, group: %s", r.consumer, r.group)
			return

		default:
			streams, err := r.rdClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    r.group,
				Consumer: r.consumer,
				Streams:  []string{r.stream, ">"},
				Count:    10,
				Block:    2 * time.Second,
			}).Result()

			if err != nil {
				if !errors.Is(err, redis.Nil) {
					log.Printf("[%s] XReadGroup error: %v", r.group, err)
					time.Sleep(time.Second)
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					callback(ctx, msg)
				}
			}
		}
	}
}

func (r *RedisStreamsConsumer) Ack(ctx context.Context, msgID string) error {
	return r.rdClient.XAck(ctx, r.stream, r.group, msgID).Err()
}
