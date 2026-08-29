package reader

import (
	rpkg "MyMessenger/pkg/redis"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type RegisterStreamsReaderToDoFunc[Event any] func(ctx context.Context, event Event) error

type readerStreams[Event any] struct {
	consumer *rpkg.RedisStreamsConsumer
	wg       sync.WaitGroup

	toDo RegisterStreamsReaderToDoFunc[Event]
}

func RegisterStreamsReader[Event any](lf fx.Lifecycle, consumer *rpkg.RedisStreamsConsumer, toDo RegisterStreamsReaderToDoFunc[Event]) {
	r := &readerStreams[Event]{toDo: toDo, consumer: consumer}

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
			log.Printf("reader: recovered from panic during event execution: %v", rec)
		}
	}()

	raw, ok := msg.Values["event"].(string)
	if !ok {
		log.Printf("readerStreams: unexpected type for event field: %T", msg.Values["event"])
		_ = r.consumer.Ack(ctx, msg.ID)
		return
	}

	var event Event
	err := json.Unmarshal([]byte(raw), &event)
	if err != nil {
		log.Printf("readerStreams: trouble with Unmarshal in callback, err: %v", err)
		_ = r.consumer.Ack(ctx, msg.ID)
		return
	}

	err = r.toDo(ctx, event)
	if err != nil {
		log.Printf("readerStreams: trouble with toDo in callback, err: %v", err)
		return
	}

	err = r.consumer.Ack(ctx, msg.ID)
	if err != nil {
		log.Printf("reader: failed to ack message %s: %v", msg.ID, err)
	}
}
