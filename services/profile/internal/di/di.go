package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	stdredis "MyMessenger/pkg/redis"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/profile/internal/config"
	"MyMessenger/services/profile/internal/infra/postgres"
	"MyMessenger/services/profile/internal/infra/postgres/migrations"
	"MyMessenger/services/profile/internal/infra/redis"
	"MyMessenger/services/profile/internal/service"
	web "MyMessenger/services/profile/internal/transport/http"

	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisClientModule,

		// Configs
		fx.Provide(
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
			config.GetProfileConfig,
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

		// Profile Service
		fx.Provide(
			stdredis.NewRedisStreamsPublisher,
			fx.Annotate(
				redis.NewRedisStreamsPublisher,
				fx.As(new(service.ProfilePublisher)),
			),
			fx.Annotate(
				postgres.NewRepo,
				fx.As(new(service.ProfileRepo)),
			),
			fx.Annotate(
				service.NewProfileService,
				fx.As(new(service.ProfileService)),
			),
		),

		// HTTP
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
