package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	stdjwt "MyMessenger/pkg/jwt"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra/postgres"
	"MyMessenger/services/auth/internal/infra/postgres/migrations"
	"MyMessenger/services/auth/internal/infra/security/argon2id"
	"MyMessenger/services/auth/internal/infra/security/jwt"
	"MyMessenger/services/auth/internal/service"
	web "MyMessenger/services/auth/internal/transport/http"
	"net/http"

	"go.uber.org/fx"
)

func provideChecker(j *stdjwt.JWTChecker) service.TokenChecker {
	return service.TokenChecker(j)
}

func GetModule() fx.Option {
	return fx.Options(

		di.JwtCheckerModule,

		// Configs
		fx.Provide(
			config.GetAuthConfig,
			config.GetHashConfig,
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
		),

		// Pool
		fx.Provide(
			repo.NewPool,
		),

		// Postgres migrations
		fx.Invoke(
			migrations.MakeMigrations,
		),

		// Auth Service
		fx.Provide(
			// Repo
			fx.Annotate(
				postgres.NewRepo,
				fx.As(new(service.AuthRepo)),
			),
			// Checker
			provideChecker,
			// Generator
			fx.Annotate(
				jwt.GetNewJwtGenerator,
				fx.As(new(service.TokenGenerator)),
			),
			// Hasher
			fx.Annotate(
				argon2id.NewArgon2idHasher,
				fx.As(new(service.AuthHasher)),
			),
			// Service
			fx.Annotate(
				service.NewJwtAuth,
				fx.As(new(service.AuthService)),
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
