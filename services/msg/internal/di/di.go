package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/msg/internal/config"
	"MyMessenger/services/msg/internal/infra/postgres"
	"MyMessenger/services/msg/internal/infra/postgres/migrations"
	"MyMessenger/services/msg/internal/infra/redis"
	"MyMessenger/services/msg/internal/service"
	"MyMessenger/services/msg/internal/transport/event"
	web "MyMessenger/services/msg/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisPublisherModule,

		// Configs
		fx.Provide(
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
			config.GetMessageConfig,
		),

		// Validator
		fx.Invoke(
			config.SetupValidator,
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
				redis.NewPublisher,
				fx.As(new(service.EventPublisher)),
			),
		),

		// Msg Service
		fx.Provide(
			fx.Annotate(
				service.NewMessageService,
				fx.As(new(service.MessageService)),
			),
		),

		// HTTP
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
