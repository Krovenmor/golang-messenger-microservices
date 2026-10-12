package ban

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"context"
	"time"
)

type Store struct {
	logins MiddlewareRepo
	tokens MiddlewareRepo
}

func NewStore(repoFab func(prefix string) MiddlewareRepo, conf *config.RedisConfig) *Store {
	return &Store{
		logins: repoFab(conf.RepoLoginPrefix),
		tokens: repoFab(conf.RepoTokenPrefix),
	}
}

func (s *Store) CheckLoginBanned(ctx context.Context, login string) bool {
	return s.logins.IsExists(ctx, domain.NormalizeLogin(login))
}

func (s *Store) CheckTokenBanned(ctx context.Context, rToken string) bool {
	return s.tokens.IsExists(ctx, utils.Hash(rToken))
}

func (s *Store) BanLogin(ctx context.Context, login string, ttl time.Duration) error {
	return s.logins.PutKey(ctx, login, ttl)
}

func (s *Store) BanTokens(ctx context.Context, rTokenHashes []string, ttl time.Duration) error {
	return s.tokens.PutKeys(ctx, rTokenHashes, ttl)
}
