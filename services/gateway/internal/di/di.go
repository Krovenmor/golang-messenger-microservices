package di

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/di"
	pkgredis "MyMessenger/pkg/redis"
	"MyMessenger/pkg/redis/redis_repo"
	"MyMessenger/services/gateway/internal/config"
	"MyMessenger/services/gateway/internal/infra/cache"
	"MyMessenger/services/gateway/internal/infra/pub"
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
	rf := redis_repo.NewRedisRepoFactory(rd, appPrefix)
	return func(prefix string) middlware.MiddlewareRepo {
		return rf.NewRedisRepo(prefix)
	}
}

func provideCacheFactory(conf *config.InfraConfig) func() middlware.MiddlewareCache {
	return func() middlware.MiddlewareCache {
		cache, err := cache.NewLruCache(conf)
		if err != nil {
			log.Fatalf("Trouble with NewLruCache, err: %q", err)
		}
		return cache
	}
}

func providePublisher(publisher *pkgredis.RedisPublisher, conf *stdconfig.RedisChannelsConfig) middlware.Publisher {
	return pub.NewRedisPublisher(publisher, conf)
}

func GetModule() fx.Option {
	return fx.Options(

		di.AuthenticatorModule,
		di.RedisPubSubModule,

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
			providePublisher,
			middlware.NewMiddleware,
		),

		// Handler
		fx.Provide(
			http.NewServeMux,
			web.NewMainHandler,
			(*web.MainHandler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			di.InvokeServer,
		),
	)
}
