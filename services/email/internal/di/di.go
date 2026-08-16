package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis"
	"MyMessenger/services/email/internal/config"
	"MyMessenger/services/email/internal/service"
	"MyMessenger/services/email/internal/transport/event"

	"go.uber.org/fx"
)

func provideSubscriber(sub *redis.RedisSubscriber) event.Subscriber {
	return event.Subscriber(sub)
}

func GetModule() fx.Option {
	return fx.Options(
		di.RedisSubscriberModule,

		fx.Provide(
			stdconfig.GetRedisChannelsConfig,
			config.GetEmailConfig,
		),

		// Email Service
		fx.Provide(
			provideSubscriber,
			fx.Annotate(
				service.NewEmailSender,
				fx.As(new(service.EmailService)),
			),
		),

		fx.Invoke(event.NewConsumer),
	)
}
