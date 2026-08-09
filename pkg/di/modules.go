package di

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/jwt"
	"MyMessenger/pkg/redis"

	"go.uber.org/fx"
)

var JwtCheckerModule = fx.Options(
	fx.Provide(
		config.GetJwtCheckerConf,
		jwt.NewJwtCheckerConf,
	),
)

var AuthenticatorModule = fx.Options(
	JwtCheckerModule,
	fx.Provide(
		jwt.NewAuthenticator,
	),
)

var RedisClientModule = fx.Option(
	fx.Provide(
		config.GetRedisConfig,
		redis.NewRedisClient,
	),
)

var RedisSubscriberModule = fx.Options(
	RedisClientModule,
	fx.Provide(
		redis.NewRedisSubscriber,
	),
)

var RedisPublisherModule = fx.Options(
	RedisClientModule,
	fx.Provide(
		redis.NewRedisPublisher,
	),
)
