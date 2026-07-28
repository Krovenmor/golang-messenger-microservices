package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra"
	"MyMessenger/services/auth/internal/infra/migrations"
	"MyMessenger/services/auth/internal/service"
	web "MyMessenger/services/auth/internal/transport/http"
	"net/http"

	"go.uber.org/fx"
)

func GetModule() fx.Option {
	return fx.Options(

		// Configs
		fx.Provide(
			config.GetAuthConfig,
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
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
				infra.NewRepo,
				fx.As(new(service.AuthRepo)),
			),
		),

		// Auth Service
		fx.Provide(
			fx.Annotate(
				service.NewAuth,
				fx.As(new(service.AuthService)),
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
