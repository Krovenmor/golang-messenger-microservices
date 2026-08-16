package redis_repo

import (
	"context"
	"fmt"
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

func (r *RedisRepo) PutKey(ctx context.Context, key string, ttl time.Duration) error {
	err := r.rd.Set(ctx, r.toKey(key), "", ttl).Err()
	if err != nil {
		return fmt.Errorf("PutKey: trouble with set, err: %w", err)
	}
	return nil
}

func (r *RedisRepo) PutVal(ctx context.Context, key, val string, ttl time.Duration) error {
	err := r.rd.Set(ctx, r.toKey(key), val, ttl).Err()
	if err != nil {
		return fmt.Errorf("PutVal: trouble with set, err: %w", err)
	}
	return nil
}

func (r *RedisRepo) PutKeys(ctx context.Context, keys []string, ttl time.Duration) error {
	if len(keys) == 0 {
		return nil
	}

	_, err := r.rd.Pipelined(ctx, func(p redis.Pipeliner) error {
		for _, key := range keys {
			p.Set(ctx, r.toKey(key), "", ttl).Err()
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("trouble with pipelined, err: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetTtl(ctx context.Context, key string) (time.Duration, bool) {
	val, err := r.rd.TTL(ctx, r.toKey(key)).Result()
	if err != nil || val == -2 {
		return 0, false
	}
	return val, true
}

func (r *RedisRepo) GetVal(ctx context.Context, key string) (string, error) {
	val, err := r.rd.Get(ctx, r.toKey(key)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisRepo) IsExists(ctx context.Context, key string) bool {
	count, err := r.rd.Exists(ctx, r.toKey(key)).Result()
	return err == nil && count > 0
}
