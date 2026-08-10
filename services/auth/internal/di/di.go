package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	stdjwt "MyMessenger/pkg/jwt"
	redispkg "MyMessenger/pkg/redis"
	rr "MyMessenger/pkg/redis/redis_repo"
	"MyMessenger/pkg/repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra/postgres"
	"MyMessenger/services/auth/internal/infra/postgres/migrations"
	"MyMessenger/services/auth/internal/infra/security/argon2id"
	"MyMessenger/services/auth/internal/infra/security/ban"
	"MyMessenger/services/auth/internal/infra/security/jwt"
	"MyMessenger/services/auth/internal/service"
	web "MyMessenger/services/auth/internal/transport/http"
	"net/http"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

const (
	appPrefix = "auth"
)

func provideChecker(j *stdjwt.JWTChecker) service.TokenChecker {
	return service.TokenChecker(j)
}

func provideRepoFabric(rd *redis.Client) func(prefix string) ban.MiddlewareRepo {
	rf := rr.NewRedisRepoFactory(rd, appPrefix)
	return func(prefix string) ban.MiddlewareRepo {
		return rf.NewRedisRepo(prefix)
	}
}

func GetModule() fx.Option {
	return fx.Options(

		di.JwtCheckerModule,
		di.RedisClientModule,

		// Configs
		fx.Provide(
			config.GetAuthConfig,
			config.GetHashConfig,
			config.GetBanConfig,
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
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

		// BanChecker
		fx.Provide(
			provideRepoFabric,
			fx.Annotate(
				redispkg.NewRedisSubscriber,
				fx.As(new(ban.Subscriber)),
			),
			fx.Annotate(
				ban.NewBanChecker,
				fx.As(new(web.BanChecker)),
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
