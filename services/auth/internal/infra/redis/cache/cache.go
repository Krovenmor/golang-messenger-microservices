package cache

import (
	rr "MyMessenger/pkg/redis/redis_repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTtlCache struct {
	repo *rr.RedisRepo
}

func NewRedisTtlCache(rf *rr.RedisFactory, conf *config.RedisConfig) *RedisTtlCache {
	return &RedisTtlCache{repo: rf.NewRedisRepo(conf.CacheTTLPrefix)}
}

func (r *RedisTtlCache) Put(ctx context.Context, email, code string, ttl time.Duration) error {
	isPut, err := r.repo.PutValNX(ctx, email, code, ttl)
	if err != nil {
		return err
	}
	if !isPut {
		return domain.ErrAlreadyExists
	}
	return nil
}

func (r *RedisTtlCache) Get(ctx context.Context, email string) (string, error) {
	val, err := r.repo.GetVal(ctx, email)
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrNotFound
	}
	return val, err
}

func (r *RedisTtlCache) Del(ctx context.Context, email string) error {
	return r.repo.DelKey(ctx, email)
}
