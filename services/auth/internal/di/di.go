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
	"MyMessenger/services/auth/internal/infra/redis/cache"
	"MyMessenger/services/auth/internal/infra/redis/pub"
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

func provideRedisRepoFactory(rd *redis.Client) *rr.RedisFactory {
	return rr.NewRedisRepoFactory(rd, appPrefix)
}

func provideRepoFabricFunc(f *rr.RedisFactory) func(prefix string) ban.MiddlewareRepo {
	return func(prefix string) ban.MiddlewareRepo {
		return f.NewRedisRepo(prefix)
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
			config.GetRedisConfig,
			stdconfig.GetRepoConfig,
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
		),

		fx.Invoke(
			config.SetupValidator,
		),

		// Pool
		fx.Provide(
			repo.NewPool,
		),

		// Postgres migrations
		fx.Invoke(
			migrations.MakeMigrations,
		),

		// Redis repo factory
		fx.Provide(
			provideRedisRepoFactory,
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
			// TTLCache
			fx.Annotate(
				cache.NewRedisTtlCache,
				fx.As(new(service.AuthTTLCache)),
			),
			// Publisher
			redispkg.NewRedisPublisher,
			fx.Annotate(
				pub.NewPublisher,
				fx.As(new(service.Publisher)),
			),
			// Service
			fx.Annotate(
				service.NewJwtAuth,
				fx.As(new(service.AuthService)),
			),
		),

		// BanChecker
		fx.Provide(
			provideRepoFabricFunc,
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
			// For Middleware
			fx.Annotate(
				stdjwt.NewAuthenticator,
				fx.As(new(web.Authenticator)),
			),
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
