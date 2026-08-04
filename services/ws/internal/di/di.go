package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis"

	web "MyMessenger/services/ws/internal/transport/http"
	"MyMessenger/services/ws/internal/transport/ws"

	"net/http"

	"go.uber.org/fx"
)

func ProvideSubscribeService(s *redis.RedisSubscriber) ws.Subscriber {
	return ws.Subscriber(s)
}

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisSubscriberModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
		),

		// Msg Broker
		fx.Provide(
			ProvideSubscribeService,
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
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
			di.InvokeServer,
		),
	)
}
