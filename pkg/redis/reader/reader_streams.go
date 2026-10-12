package reader

import (
	"MyMessenger/pkg/logger"
	rpkg "MyMessenger/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RegisterStreamsReaderToDoFunc[Event any] func(ctx context.Context, event Event) error

type readerStreams[Event any] struct {
	consumer *rpkg.RedisStreamsConsumer
	wg       sync.WaitGroup
	log      *logger.Logger

	toDo RegisterStreamsReaderToDoFunc[Event]
}

func RegisterStreamsReader[Event any](lf fx.Lifecycle, log *logger.Logger, consumer *rpkg.RedisStreamsConsumer, toDo RegisterStreamsReaderToDoFunc[Event]) {
	r := &readerStreams[Event]{toDo: toDo, consumer: consumer, log: log}

	ctxC, ctxCancel := context.WithCancel(context.Background())

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := consumer.InitGroup(ctx)
			if err != nil {
				return err
			}

			r.wg.Add(1)
			go func() {
				defer r.wg.Done()
				consumer.Start(ctxC, r.callback)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctxCancel()

			r.wg.Wait()
			return nil
		},
	})
}

func (r *readerStreams[Event]) callback(ctx context.Context, msg redis.XMessage) {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("recovered from panic during event execution", zap.Any("panic", rec), zap.Stack("stack"))
		}
	}()

	raw, ok := msg.Values["event"].(string)
	if !ok {
		r.log.Error("unexpected type for event field", zap.String("type", fmt.Sprintf("%T", msg.Values["event"])), zap.String("message_id", msg.ID))
		_ = r.consumer.Ack(ctx, msg.ID)
		return
	}

	var event Event
	err := json.Unmarshal([]byte(raw), &event)
	if err != nil {
		r.log.Error("failed to unmarshal event", zap.String("message_id", msg.ID), zap.Error(err))
		_ = r.consumer.Ack(ctx, msg.ID)
		return
	}

	err = r.toDo(ctx, event)
	if err != nil {
		r.log.Ctx(ctx).Error("failed to handle event", zap.String("message_id", msg.ID), zap.Error(err))
		return
	}

	err = r.consumer.Ack(ctx, msg.ID)
	if err != nil {
		r.log.Ctx(ctx).Error("failed to ack message", zap.String("message_id", msg.ID), zap.Error(err))
	}
}
