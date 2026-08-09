package middlware

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type MiddlewareRepo interface {
	Put(key any, expAt time.Time)
	Get(key any) (time.Time, bool)
	IsExists(key any) bool
}

type MiddlewareCache interface {
	Put(key string, limiter *rate.Limiter)
	Get(key string) (*rate.Limiter, bool)
}

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}
