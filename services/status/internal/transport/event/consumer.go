package event

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/services/status/internal/config"
	"MyMessenger/services/status/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"go.uber.org/fx"
)

type Subscriber interface {
	PatternSubscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

type Consumer struct {
	stService  service.StatusService
	sub        Subscriber
	subChannel string
	wg         sync.WaitGroup
}

func NewConsumer(stService service.StatusService, sub Subscriber, conf *config.SubInfoConfig) *Consumer {
	return &Consumer{
		stService:  stService,
		sub:        sub,
		subChannel: fmt.Sprintf(conf.SubPattern, "*"),
	}
}

func (c *Consumer) startConsuming(ctx context.Context) {
	defer c.wg.Done()

	ch, stop, err := c.sub.PatternSubscribe(ctx, c.subChannel)
	if err != nil {
		log.Printf("Trouble with c.sub.PatternSubscribe(ctx, %q), err %q", c.subChannel, err.Error())
		return
	}
	defer stop()

	log.Printf("Starting consuming on channel: %q", c.subChannel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("startConsuming: context ended")
			return
		case data, ok := <-ch:
			if !ok {
				log.Printf("startConsuming: data stream ended, channel closed")
				return
			}

			var event broker.StatusEvent
			err := json.Unmarshal(data, &event)
			if err != nil {
				log.Printf("startConsuming: Bad JSON from broker, data: %q", data)
				continue
			}

			log.Printf("startConsuming: New status event, event: %v", event)

			err = c.stService.SaveStatus(ctx, event)
			if err != nil {
				log.Printf("startConsuming: trouble with saving event: %q", err.Error())
			}
		}
	}
}

func (c *Consumer) Consume(lf fx.Lifecycle) {
	bCtx, cancel := context.WithCancel(context.Background())
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			c.wg.Add(1)
			go c.startConsuming(bCtx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			c.wg.Wait()
			return nil
		},
	})
}
