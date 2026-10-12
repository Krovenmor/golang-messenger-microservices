package reader

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/logger"
	pkgredis "MyMessenger/pkg/redis"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CallbackReader interface {
	GetServiceName() string

	OnNewProfileEvent(ctx context.Context, event broker.ProfilePayload) error
}

type RedisProfileReader struct {
	reader  CallbackReader
	infoStr string
	log     *logger.Logger
}

func NewRedisProfileReader(lf fx.Lifecycle, log *logger.Logger,
	rdClient *redis.Client, reader CallbackReader,
	conf *config.RedisChannelsConfig) *RedisProfileReader {

	r := &RedisProfileReader{
		reader: reader,
		log:    log,
	}

	serviceName := reader.GetServiceName()
	consumerName := pkgredis.GenUniqueConsumerName(serviceName)
	groupName := serviceName + "-group"

	r.infoStr = fmt.Sprintf("info:stream:%q:group:%q:consumer:%q", conf.ProfileStream, groupName, consumerName)

	consumer := pkgredis.NewRedisStreamsConsumer(rdClient, log, conf.ProfileStream, groupName, consumerName)
	RegisterStreamsReader(lf, log, consumer, r.onEvent)

	return r
}

func (r *RedisProfileReader) toEventErr(err error) error {
	return fmt.Errorf("profile event: %w: %s", err, r.infoStr)
}

func (r *RedisProfileReader) onEvent(ctx context.Context, event broker.ProfileEventDTO) error {
	switch event.Type {
	case broker.NewProfileEvent:
		err := r.reader.OnNewProfileEvent(ctx, event.Payload)
		if err != nil {
			return r.toEventErr(err)
		}
	default:
		r.log.Ctx(ctx).Warn("unknown profile event", zap.String("type", string(event.Type)))
		return nil
	}
	r.log.Ctx(ctx).Debug("profile event handled", zap.String("type", string(event.Type)))
	return nil
}
