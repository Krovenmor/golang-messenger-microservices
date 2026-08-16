package cache

import (
	rr "MyMessenger/pkg/redis/redis_repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/service"
	"context"
	"time"
)

type RedisTtlCache struct {
	repo *rr.RedisRepo
}

func NewRedisTtlCache(rf *rr.RedisFactory, conf *config.RedisConfig) *RedisTtlCache {
	return &RedisTtlCache{repo: rf.NewRedisRepo(conf.CacheTTLPrefix)}
}

func (r *RedisTtlCache) Put(ctx context.Context, email, code string, ttl time.Duration) error {
	if r.repo.IsExists(ctx, email) {
		return service.ErrAlreadyExists
	}
	return r.repo.PutVal(ctx, email, code, ttl)
}

func (r *RedisTtlCache) Get(ctx context.Context, email string) (string, error) {
	return r.repo.GetVal(ctx, email)
}

func (r *RedisTtlCache) Del(ctx context.Context, email string) error {
	return r.repo.DelKey(ctx, email)
}
