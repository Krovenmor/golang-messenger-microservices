package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	pkgredis "MyMessenger/pkg/redis"
	"MyMessenger/services/gateway/internal/config"
	"MyMessenger/services/gateway/internal/infra"
	"MyMessenger/services/gateway/internal/middlware"
	web "MyMessenger/services/gateway/internal/transport/http"
	"log"

	"net/http"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

const (
	appPrefix = "gtw"
)

func provideSub(sub *pkgredis.RedisSubscriber) middlware.Subscriber {
	return middlware.Subscriber(sub)
}

func provideRepoFactory(rd *redis.Client) func(prefix string) middlware.MiddlewareRepo {
	rf := infra.NewRedisRepoFactory(rd, appPrefix)
	return func(prefix string) middlware.MiddlewareRepo {
		return rf.NewRedisRepo(prefix)
	}
}

func provideCacheFactory(conf *config.InfraConfig) func() middlware.MiddlewareCache {
	return func() middlware.MiddlewareCache {
		cache, err := infra.NewLruCache(conf)
		if err != nil {
			log.Fatalf("Trouble with NewLruCache, err: %q", err)
		}
		return cache
	}
}

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisSubscriberModule,

		// Configs
		fx.Provide(
			stdconfig.GetServConfig,
			stdconfig.GetRedisChannelsConfig,
			config.GetMainHandlerConfig,
			config.GetMiddlewareConfig,
			config.GetInfraConfig,
		),

		// Middleware
		fx.Provide(
			provideSub,
			provideCacheFactory,
			provideRepoFactory,
			middlware.NewMiddleware,
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
		),

		// Handler
		fx.Provide(
			web.NewMainHandler,
			(*web.MainHandler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
		),
	)
}
