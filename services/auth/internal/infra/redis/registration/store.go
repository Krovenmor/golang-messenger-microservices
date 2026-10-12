package registration

import (
	rr "MyMessenger/pkg/redis/redis_repo"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type record struct {
	Login        string `json:"login"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	EmailTaken   bool   `json:"emailTaken"`
}

type Store struct {
	registrations *rr.RedisRepo
	codes         *rr.RedisRepo
	cooldowns     *rr.RedisRepo
}

func NewStore(rf *rr.RedisFactory, conf *config.RedisConfig) *Store {
	return &Store{
		registrations: rf.NewRedisRepo(conf.RegistrationPrefix),
		codes:         rf.NewRedisRepo(conf.RegistrationCodePrefix),
		cooldowns:     rf.NewRedisRepo(conf.CodeCooldownPrefix),
	}
}

func (s *Store) Save(ctx context.Context, reg *domain.Registration, ttl time.Duration) error {
	data, err := json.Marshal(record{
		Login:        reg.Login,
		Email:        reg.Email,
		PasswordHash: reg.PasswordHash,
		EmailTaken:   reg.EmailTaken,
	})
	if err != nil {
		return fmt.Errorf("marshal registration: %w", err)
	}
	return s.registrations.PutVal(ctx, reg.ID.String(), string(data), ttl)
}

func (s *Store) Get(ctx context.Context, id uuid.UUID) (*domain.Registration, error) {
	data, err := s.registrations.GetVal(ctx, id.String())
	if errors.Is(err, redis.Nil) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var rec record
	err = json.Unmarshal([]byte(data), &rec)
	if err != nil {
		return nil, fmt.Errorf("unmarshal registration: %w", err)
	}
	return &domain.Registration{
		ID:           id,
		Login:        rec.Login,
		Email:        rec.Email,
		PasswordHash: rec.PasswordHash,
		EmailTaken:   rec.EmailTaken,
	}, nil
}

func (s *Store) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.Join(
		s.registrations.DelKey(ctx, id.String()),
		s.codes.DelKey(ctx, id.String()),
	)
}

func (s *Store) SaveCode(ctx context.Context, id uuid.UUID, code string, ttl time.Duration) error {
	return s.codes.PutVal(ctx, id.String(), code, ttl)
}

func (s *Store) GetCode(ctx context.Context, id uuid.UUID) (string, error) {
	code, err := s.codes.GetVal(ctx, id.String())
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrNotFound
	}
	return code, err
}

func (s *Store) StartCooldown(ctx context.Context, email string, ttl time.Duration) error {
	isStarted, err := s.cooldowns.PutValNX(ctx, email, "", ttl)
	if err != nil {
		return err
	}
	if !isStarted {
		return domain.ErrAlreadyExists
	}
	return nil
}

func (s *Store) StopCooldown(ctx context.Context, email string) error {
	return s.cooldowns.DelKey(ctx, email)
}
