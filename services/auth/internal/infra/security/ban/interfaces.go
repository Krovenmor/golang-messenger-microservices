package ban

import (
	"context"
	"time"
)

type MiddlewareRepo interface {
	Put(ctx context.Context, key string, ttl time.Duration) error
	PutKeys(ctx context.Context, key []string, ttl time.Duration) error
	Get(ctx context.Context, key string) (time.Duration, bool)
	IsExists(ctx context.Context, key string) bool
}

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}
