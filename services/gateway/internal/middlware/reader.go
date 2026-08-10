package middlware

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/redis/reader"
	"context"
	"log"
	"time"

	"go.uber.org/fx"
)

func registerReader(lf fx.Lifecycle, sub Subscriber, pub Publisher, channel string, repo MiddlewareRepo, calcTime func(reason BanReason) time.Duration) {
	reader.RegisterReader(lf, sub, channel, func(ctx context.Context, event broker.BanRequestDTO) {
		if event.Type != broker.BanRequest {
			log.Printf("reader: Trouble with BanEventDTO, type want: %q, got: %q", broker.BanRequest, event.Type)
			return
		}
		reason := FromBroker(event.Payload.Reason)
		ttl := calcTime(reason)
		err := repo.Put(ctx, event.Payload.UserId.String(), ttl)
		if err != nil {
			log.Printf("reader: trouble with saving ban: %v, err: %q", event, err)
			return
		}
		log.Printf("reader: new ban event registered: %v", event)
		err = pub.Pub(ctx, broker.BanEventPayload{
			UserId: event.Payload.UserId,
			Ttl:    ttl,
		})
		if err != nil {
			log.Printf("reader: trouble with publicate ban: %v, err: %q", event, err)
			return
		}
		log.Printf("reader: ban publicated")
	})
}
