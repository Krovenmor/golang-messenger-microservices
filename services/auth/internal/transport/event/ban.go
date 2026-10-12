package event

import (
	"MyMessenger/pkg/broker"
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/logger"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/services/auth/internal/service"
	"context"

	"go.uber.org/fx"
)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

func RegisterBanConsumer(lf fx.Lifecycle, log *logger.Logger, sub Subscriber, auth service.AuthService, conf *stdconfig.RedisChannelsConfig) {
	reader.RegisterReader(lf, log, sub, conf.UserBanEventChannel, func(ctx context.Context, event broker.BanEventDTO) {
		_ = auth.BanUser(ctx, event.Payload.UserId, event.Payload.Ttl)
	})
}
