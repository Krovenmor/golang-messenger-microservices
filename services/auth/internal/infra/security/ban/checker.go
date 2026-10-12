package ban

import (
	"MyMessenger/pkg/broker"
	stdconf "MyMessenger/pkg/config"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/pkg/utils"

	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/service"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

type BanChecker struct {
	repoLogins MiddlewareRepo
	repoTokens MiddlewareRepo

	auth service.AuthService
}

func NewBanChecker(lf fx.Lifecycle, repoFab func(prefix string) MiddlewareRepo, sub Subscriber, auth service.AuthService,
	banConf *config.RedisConfig, redisConf *stdconf.RedisChannelsConfig) *BanChecker {
	checker := &BanChecker{
		repoLogins: repoFab(banConf.RepoLoginPrefix),
		repoTokens: repoFab(banConf.RepoTokenPrefix),
		auth:       auth,
	}
	reader.RegisterReader(lf, sub, redisConf.UserBanEventChannel, func(ctx context.Context, event broker.BanEventDTO) {
		checker.banUser(ctx, event.Payload.UserId, event.Payload.Ttl)
	})
	return checker
}

func (c *BanChecker) CheckLoginBanned(ctx context.Context, login string) bool {
	return c.repoLogins.IsExists(ctx, domain.NormalizeLogin(login))
}

func (c *BanChecker) CheckTokenBanned(ctx context.Context, rToken string) bool {
	return c.repoTokens.IsExists(ctx, utils.Hash(rToken))
}

func (c *BanChecker) banUser(ctx context.Context, userId uuid.UUID, ttl time.Duration) {
	info, err := c.auth.GetUserTokens(ctx, userId)
	if errors.Is(err, domain.ErrNotFound) {
		c.banLoginOnly(ctx, userId, ttl)
		return
	}
	if err != nil {
		log.Printf("banUser: trouble with GetUserTokens(%q), err: %q", userId, err)
		return
	}
	err = c.repoLogins.PutKey(ctx, info.Login, ttl)
	if err != nil {
		log.Printf("banUser: Trouble with repoLogins.PutKey(ctx, %v, %d)", userId, ttl)
	}
	err = c.repoTokens.PutKeys(ctx, info.RTokens, ttl)
	if err != nil {
		log.Printf("banUser: Trouble with repoTokens.PutKeys(ctx, %v, %d)", info.RTokens, ttl)
	}
}

func (c *BanChecker) banLoginOnly(ctx context.Context, userId uuid.UUID, ttl time.Duration) {
	info, err := c.auth.GetUserInfo(ctx, userId)
	if err != nil {
		log.Printf("banLoginOnly: trouble with GetUserInfo(%q), err: %q", userId, err)
		return
	}
	err = c.repoLogins.PutKey(ctx, info.Login, ttl)
	if err != nil {
		log.Printf("banLoginOnly: Trouble with repoLogins.PutKey(ctx, %v, %d)", userId, ttl)
	}
}
