package cache

import (
	"MyMessenger/services/gateway/internal/config"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/time/rate"
)

type LruCache struct {
	l *lru.TwoQueueCache[string, *rate.Limiter]
}

func NewLruCache(conf *config.InfraConfig) (*LruCache, error) {
	cache, err := lru.New2Q[string, *rate.Limiter](conf.CacheSize)
	if err != nil {
		return nil, err
	}

	return &LruCache{
		l: cache,
	}, nil
}

func (c *LruCache) Put(key string, limiter *rate.Limiter) {
	c.l.Add(key, limiter)
}

func (c *LruCache) Get(key string) (*rate.Limiter, bool) {
	return c.l.Get(key)
}
