package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis"

	"MyMessenger/services/ws/internal/client/msg"
	"MyMessenger/services/ws/internal/config"
	"MyMessenger/services/ws/internal/infra"
	"MyMessenger/services/ws/internal/service"
	web "MyMessenger/services/ws/internal/transport/http"
	"MyMessenger/services/ws/internal/transport/ws"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisClientModule,

		fx.Provide(
			redis.NewRedisPublisher,
			redis.NewRedisSubscriber,
		),

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
			config.GetWsConfig,
			config.GetMsgClientConfig,
		),

		// Msg Broker
		fx.Provide(
			fx.Annotate(
				infra.NewPublisher,
				fx.As(new(service.Publisher)),
			),
			fx.Annotate(
				infra.NewSubscriber,
				fx.As(new(service.Subscriber)),
			),
		),

		// Msg Client
		fx.Provide(
			fx.Annotate(
				msg.NewMsgClient,
				fx.As(new(service.MessageClient)),
			),
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
		),

		// Service
		fx.Provide(
			fx.Annotate(
				service.NewWsService,
				fx.As(new(service.WsService)),
			),
		),

		// Handlers
		fx.Provide(
			ws.NewWSHandler,
			web.NewMainHandler,
		),

		fx.Provide(
			(*web.MainHandler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			service.ComputeJsonResponses,
			di.InvokeServer,
		),
	)
}
