package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/services/media/internal/config"
	"MyMessenger/services/media/internal/infra/postgres"
	"MyMessenger/services/media/internal/infra/postgres/migrations"
	"MyMessenger/services/media/internal/service"
	web "MyMessenger/services/media/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.PgxpoolModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			config.GetMediaSaverConfig,
			config.GetHandlerConfig,
		),

		// Repo
		fx.Invoke(
			migrations.MakeMigrations,
		),

		// Service
		fx.Provide(
			fx.Annotate(
				postgres.NewPostgresRepo,
				fx.As(new(service.MediaRepo)),
			),
			fx.Annotate(
				service.NewMediaSaver,
				fx.As(new(service.MediaService)),
			),
		),

		// Handler
		fx.Provide(
			http.NewServeMux,
			web.NewHandler,
			(*web.Handler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
		),
	)
}
