package ban

import (
	"context"
	"time"
)

type MiddlewareRepo interface {
	PutKey(ctx context.Context, key string, ttl time.Duration) error
	PutKeys(ctx context.Context, key []string, ttl time.Duration) error
	IsExists(ctx context.Context, key string) bool
}

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}
