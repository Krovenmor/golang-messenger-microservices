package middlware

import (
	"MyMessenger/pkg/broker"
	"context"
	"encoding/json"
	"log"
	"time"

	"go.uber.org/fx"
)

type reader struct {
	banChannel <-chan []byte
	banClose   func()

	repo    MiddlewareRepo
	channel string
}

func registerReader(lf fx.Lifecycle, sub Subscriber, channel string, repo MiddlewareRepo, calcTime func(reason BanReason) time.Duration) {
	r := &reader{repo: repo, channel: channel}

	ctxC, ctxCancel := context.WithCancel(context.Background())

	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ch, cl, err := sub.Subscribe(ctxC, channel)
			if err != nil {
				return err
			}
			r.banChannel = ch
			r.banClose = cl
			go r.startReader(ctxC, calcTime)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if r.banClose != nil {
				r.banClose()
			}
			ctxCancel()
			return nil
		},
	})
}

func (r *reader) startReader(ctx context.Context, calcTime func(reason BanReason) time.Duration) {
	log.Printf("Reader started on %q", r.channel)
	defer log.Printf("Reader stopped on %q", r.channel)

	for {
		select {
		case msg, ok := <-r.banChannel:
			if !ok {
				return
			}
			var event broker.BanEventDTO
			err := json.Unmarshal(msg, &event)
			if err != nil {
				log.Printf("reader: Trouble with unmarshaling BanEventDTO, err: %q", err)
				continue
			}
			if event.Type != broker.BanEvent {
				log.Printf("reader: Trouble with BanEventDTO, type want: %q, got: %q", broker.BanEvent, event.Type)
				continue
			}
			reason := FromBroker(event.Payload.Reason)
			r.repo.Put(ctx, event.Payload.UserId.String(), calcTime(reason))
			log.Printf("reader: new ban event registered: %v", event)

		case <-ctx.Done():
			return
		}
	}
}
