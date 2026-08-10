package reader

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/fx"
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

type reader[Event any] struct {
	ch      <-chan []byte
	closeCh func()
	channel string
	wg      sync.WaitGroup

	toDo func(ctx context.Context, event Event)
}

func RegisterReader[Event any](lf fx.Lifecycle, sub Subscriber, channel string, toDo func(ctx context.Context, event Event)) {
	r := &reader[Event]{channel: channel, toDo: toDo}

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
	log.Printf("Reader started on %q", r.channel)
	defer log.Printf("Reader stopped on %q", r.channel)

	for {
		select {
		case msg, ok := <-r.ch:
			if !ok {
				return
			}
			var event Event
			err := json.Unmarshal(msg, &event)
			if err != nil {
				log.Printf("reader [%s]: Trouble with unmarshaling Event, err: %q", r.channel, err)
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
			log.Printf("reader [%s]: recovered from panic during event execution: %v", r.channel, rec)
		}
	}()

	r.toDo(ctx, event)
}
