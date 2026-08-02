package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"

	"MyMessenger/services/ws/internal/infra/redis"

	web "MyMessenger/services/ws/internal/transport/http"
	"MyMessenger/services/ws/internal/transport/ws"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			stdconfig.GetRedisConfig,
		),

		// Msg Broker
		fx.Provide(
			fx.Annotate(
				redis.NewRedisSubscriber,
				fx.As(new(ws.Subscriber)),
			),
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
