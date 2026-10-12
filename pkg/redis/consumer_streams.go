package redis

import (
	"MyMessenger/pkg/logger"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CallbackFunc = func(ctx context.Context, msg redis.XMessage)

type RedisStreamsConsumer struct {
	rdClient *redis.Client
	log      *logger.Logger

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

func NewRedisStreamsConsumer(rdClient *redis.Client, log *logger.Logger, stream, group, consumer string) *RedisStreamsConsumer {
	return &RedisStreamsConsumer{
		rdClient: rdClient,
		log:      log,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (r *RedisStreamsConsumer) InitGroup(ctx context.Context) error {
	err := r.rdClient.XGroupCreateMkStream(ctx, r.stream, r.group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("create consumer group %q: %w", r.group, err)
	}
	return nil
}

func (r *RedisStreamsConsumer) Start(ctx context.Context, callback CallbackFunc) {
	for {
		select {
		case <-ctx.Done():
			r.log.Info("streams consumer stopped", zap.String("consumer", r.consumer), zap.String("group", r.group))
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
					r.log.Error("streams read failed", zap.String("group", r.group), zap.Error(err))
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
