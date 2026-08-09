package middlware

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type MiddlewareRepo interface {
	Put(ctx context.Context, key string, ttl time.Duration)
	Get(ctx context.Context, key string) (time.Duration, bool)
	IsExists(ctx context.Context, key string) bool
}

type MiddlewareCache interface {
	Put(key string, limiter *rate.Limiter)
	Get(key string) (*rate.Limiter, bool)
}

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}
