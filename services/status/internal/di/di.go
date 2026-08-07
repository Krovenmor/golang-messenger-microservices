package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis"

	config "MyMessenger/services/status/internal/config"
	"MyMessenger/services/status/internal/infra"
	"MyMessenger/services/status/internal/service"
	"MyMessenger/services/status/internal/transport/event"
	web "MyMessenger/services/status/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func ProvideSubscribeService(s *redis.RedisSubscriber) event.Subscriber {
	return event.Subscriber(s)
}

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisSubscriberModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
			config.GetServiceConfig,
		),

		// Repo
		fx.Provide(
			fx.Annotate(
				infra.NewRedisRepo,
				fx.As(new(service.StatusRepo)),
			),
		),

		// Msg Broker
		fx.Provide(
			ProvideSubscribeService,
		),

		// Domain
		fx.Provide(
			fx.Annotate(
				service.NewStatusService,
				fx.As(new(service.StatusService)),
			),
		),

		// Http startup
		fx.Provide(
			http.NewServeMux,
			web.NewHandler,
			(*web.Handler).RegisterRoutes,
		),

		// Consumer startup
		fx.Provide(
			event.NewConsumer,
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
			(*event.Consumer).Consume,
		),
	)
}
