package event

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/services/email/internal/service"
	"context"
	"log"

	"go.uber.org/fx"
)

const (
	codeSize = 6
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

type Consumer struct {
	es service.EmailService
}

func NewConsumer(lf fx.Lifecycle, sub Subscriber, es service.EmailService, conf *config.RedisChannelsConfig) *Consumer {
	consumer := Consumer{es: es}
	reader.RegisterReader(lf, sub, conf.UserVerificationChannel, consumer.newEmailEvent)
	return &consumer
}

func (c *Consumer) newEmailEvent(ctx context.Context, event broker.EmailVerificationDTO) {
	if event.Type != broker.EmailVerificationType {
		log.Printf("newEmailEvent: Wrong event type: %q", event.Type)
		return
	}
	c.es.SendVerificationCode(ctx, event.Payload.Email, event.Payload.Code)
}
