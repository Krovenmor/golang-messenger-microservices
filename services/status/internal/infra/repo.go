package infra

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/pkg/config"
	"MyMessenger/services/status/internal/service"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	keyPatternPrefix = "status:"
	saveParts        = 2
)

type RedisRepo struct {
	rd      *redis.Client
	pattern string
}

func NewRedisRepo(rdClient *redis.Client, conf *config.RedisChannelsConfig) *RedisRepo {
	return &RedisRepo{rd: rdClient, pattern: keyPatternPrefix + conf.UserStatusPattern}
}

func (r *RedisRepo) ToKey(userId string) string {
	return fmt.Sprintf(r.pattern, userId)
}

func (r *RedisRepo) ToSave(s *broker.StatusPayload) string {
	return fmt.Sprintf("%d:%d", s.Status, s.EventTime)
}

func (r *RedisRepo) FromSaved(s string) *service.UserStatus {
	sl := strings.Split(s, ":")
	if len(sl) != saveParts {
		return nil
	}
	status, err := strconv.ParseInt(sl[0], 10, 64)
	if err != nil {
		return nil
	}
	lastSeen, err := strconv.ParseInt(sl[1], 10, 64)
	if err != nil {
		return nil
	}
	return &service.UserStatus{
		Status:   broker.Status(status),
		LastSeen: int64(lastSeen),
	}
}

func (r *RedisRepo) GetStatus(ctx context.Context, userId string) (*service.UserStatus, error) {
	toParse, err := r.rd.Get(ctx, r.ToKey(userId)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get status for user %s: %w", userId, err)
	}
	status := r.FromSaved(toParse)
	if status == nil {
		return nil, fmt.Errorf("Trouble with parsing")
	}
	return status, nil
}

func (r *RedisRepo) SaveStatus(ctx context.Context, status broker.StatusPayload, ttl time.Duration) error {
	err := r.rd.Set(ctx, r.ToKey(status.UserId), r.ToSave(&status), ttl).Err()
	if err != nil {
		return fmt.Errorf("Trouble with saving, err: %w", err)
	}
	return nil
}
