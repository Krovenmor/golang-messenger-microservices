package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/services/media/internal/config"
	"MyMessenger/services/media/internal/infra/postgres"
	"MyMessenger/services/media/internal/infra/postgres/migrations"
	"MyMessenger/services/media/internal/service"
	"MyMessenger/services/media/internal/transport/event"
	web "MyMessenger/services/media/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.PgxpoolModule,
		di.RedisClientModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
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

		// Event
		fx.Provide(
			fx.Annotate(
				event.NewProfileReader,
				fx.As(new(reader.CallbackReader)),
			),
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
			reader.NewRedisProfileReader,
		),
	)
}
