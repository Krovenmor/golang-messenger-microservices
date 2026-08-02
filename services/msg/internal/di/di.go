package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/msg/internal/config"
	"MyMessenger/services/msg/internal/infra/postgres"
	"MyMessenger/services/msg/internal/infra/postgres/migrations"
	"MyMessenger/services/msg/internal/infra/redis"
	"MyMessenger/services/msg/internal/service"
	web "MyMessenger/services/msg/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,

		// Configs
		fx.Provide(
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
			stdconfig.GetRedisConfig,
			config.GetMessageConfig,
		),

		// Pool
		fx.Provide(
			repo.NewPool,
		),

		// Make migrations
		fx.Invoke(
			migrations.MakeMigrations,
		),

		// Repo Service
		fx.Provide(
			fx.Annotate(
				postgres.NewRepo,
				fx.As(new(service.MessageRepo)),
			),
		),

		// Msg Broker
		fx.Provide(
			fx.Annotate(
				redis.NewRedisPublisher,
				fx.As(new(service.EventPublisher)),
			),
		),

		// Msg Service
		fx.Provide(
			fx.Annotate(
				service.NewMessageServiceImpl,
				fx.As(new(service.MessageService)),
			),
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
		),

		fx.Provide(
			web.NewHandler,
		),

		fx.Provide(
			(*web.Handler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
		),
	)
}
