package reader

import (
	"MyMessenger/pkg/logger"
	"context"
	"encoding/json"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

type reader[Event any] struct {
	ch      <-chan []byte
	closeCh func()
	channel string
	wg      sync.WaitGroup
	log     *zap.Logger

	toDo func(ctx context.Context, event Event)
}

func RegisterReader[Event any](lf fx.Lifecycle, log *logger.Logger, sub Subscriber, channel string, toDo func(ctx context.Context, event Event)) {
	r := &reader[Event]{channel: channel, toDo: toDo, log: log.With(zap.String("channel", channel))}

	ctxC, ctxCancel := context.WithCancel(context.Background())

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ch, cl, err := sub.Subscribe(ctxC, channel)
			if err != nil {
				return err
			}
			r.ch = ch
			r.closeCh = cl

			r.wg.Add(1)
			go func() {
				defer r.wg.Done()
				r.startReader(ctxC)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if r.closeCh != nil {
				r.closeCh()
			}
			ctxCancel()

			r.wg.Wait()
			return nil
		},
	})
}

func (r *reader[Event]) startReader(ctx context.Context) {
	r.log.Info("reader started")
	defer r.log.Info("reader stopped")

	for {
		select {
		case msg, ok := <-r.ch:
			if !ok {
				return
			}
			var event Event
			err := json.Unmarshal(msg, &event)
			if err != nil {
				r.log.Error("failed to unmarshal event", zap.Error(err))
				continue
			}
			r.toExec(ctx, event)

		case <-ctx.Done():
			return
		}
	}
}

func (r *reader[Event]) toExec(ctx context.Context, event Event) {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("recovered from panic during event execution", zap.Any("panic", rec), zap.Stack("stack"))
		}
	}()

	r.toDo(ctx, event)
}
