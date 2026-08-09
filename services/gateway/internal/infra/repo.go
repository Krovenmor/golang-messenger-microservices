package infra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFactory struct {
	rd        *redis.Client
	appPrefix string
}

type RedisRepo struct {
	rd     *redis.Client
	prefix string
}

func NewRedisRepoFactory(rd *redis.Client, appPrefix string) *RedisFactory {
	return &RedisFactory{
		rd:        rd,
		appPrefix: appPrefix,
	}
}

func (f *RedisFactory) NewRedisRepo(modulePrefix string) *RedisRepo {
	fullPrefix := fmt.Sprintf("%s:%s", f.appPrefix, modulePrefix)
	return &RedisRepo{
		rd:     f.rd,
		prefix: fullPrefix,
	}
}

func (r *RedisRepo) toKey(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

func (r *RedisRepo) Put(ctx context.Context, key string, ttl time.Duration) {
	err := r.rd.Set(ctx, r.toKey(key), "b", ttl).Err()
	if err != nil {
		log.Printf("Put: Trouble with set, err: %q", err)
	}
}

func (r *RedisRepo) Get(ctx context.Context, key string) (time.Duration, bool) {
	val, err := r.rd.TTL(ctx, r.toKey(key)).Result()
	if err != nil || val == -2 {
		return 0, false
	}
	return val, true
}

func (r *RedisRepo) IsExists(ctx context.Context, key string) bool {
	count, err := r.rd.Exists(ctx, r.toKey(key)).Result()
	return err == nil && count > 0
}
