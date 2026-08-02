package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/services/gateway/internal/config"
	"MyMessenger/services/gateway/internal/middlware"
	web "MyMessenger/services/gateway/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			config.GetMainHandlerConfig,
		),

		// Checkers
		fx.Provide(
			middlware.NewLimitChecker,
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
		),

		// Handler
		fx.Provide(
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
