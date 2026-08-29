package reader

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	pkgredis "MyMessenger/pkg/redis"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type CallbackReader interface {
	GetServiceName() string

	OnNewProfileEvent(ctx context.Context, event broker.ProfilePayload) error
}

type RedisProfileReader struct {
	reader  CallbackReader
	infoStr string
}

func NewRedisProfileReader(lf fx.Lifecycle,
	rdClient *redis.Client, reader CallbackReader,
	conf *config.RedisChannelsConfig) *RedisProfileReader {

	r := &RedisProfileReader{
		reader: reader,
	}

	serviceName := reader.GetServiceName()
	consumerName := pkgredis.GenUniqueConsumerName(serviceName)
	groupName := serviceName + "-group"

	r.infoStr = fmt.Sprintf("info:stream:%q:group:%q:consumer:%q", conf.ProfileStream, groupName, consumerName)

	consumer := pkgredis.NewRedisStreamsConsumer(rdClient, conf.ProfileStream, groupName, consumerName)
	RegisterStreamsReader(lf, consumer, r.onEvent)

	return r
}

func (r *RedisProfileReader) toEventErr(err error) error {
	err = fmt.Errorf("redisReader.onEvent(), err:%w:%s", err, r.infoStr)
	log.Printf("%q", err)
	return err
}

func (r *RedisProfileReader) onEvent(ctx context.Context, event broker.ProfileEventDTO) error {
	switch event.Type {
	case broker.NewProfileEvent:
		err := r.reader.OnNewProfileEvent(ctx, event.Payload)
		if err != nil {
			return r.toEventErr(err)
		}
	default:
		log.Printf("redisReader.onEvent(): unknown event=%q", event.Type)
	}
	return nil
}
