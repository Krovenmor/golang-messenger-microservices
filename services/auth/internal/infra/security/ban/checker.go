package ban

import (
	"MyMessenger/pkg/broker"
	stdconf "MyMessenger/pkg/config"
	"MyMessenger/pkg/logger"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/service"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BanChecker struct {
	repoLogins MiddlewareRepo
	repoTokens MiddlewareRepo

	auth service.AuthService

	log *logger.Logger
}

func NewBanChecker(lf fx.Lifecycle, log *logger.Logger, repoFab func(prefix string) MiddlewareRepo, sub Subscriber, auth service.AuthService,
	banConf *config.RedisConfig, redisConf *stdconf.RedisChannelsConfig) *BanChecker {
	checker := &BanChecker{
		repoLogins: repoFab(banConf.RepoLoginPrefix),
		repoTokens: repoFab(banConf.RepoTokenPrefix),
		auth:       auth,
		log:        log,
	}
	reader.RegisterReader(lf, log, sub, redisConf.UserBanEventChannel, func(ctx context.Context, event broker.BanEventDTO) {
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
	log := c.log.Ctx(ctx).With(zap.Stringer("user_id", userId), zap.Duration("ttl", ttl))

	info, err := c.auth.GetUserTokens(ctx, userId)
	if errors.Is(err, domain.ErrNotFound) {
		c.banLoginOnly(ctx, log, userId, ttl)
		return
	}
	if err != nil {
		log.Error("failed to get user tokens for ban", zap.Error(err))
		return
	}
	err = c.repoLogins.PutKey(ctx, info.Login, ttl)
	if err != nil {
		log.Error("failed to ban login", zap.Error(err))
	}
	err = c.repoTokens.PutKeys(ctx, info.RTokens, ttl)
	if err != nil {
		log.Error("failed to ban refresh tokens", zap.Int("tokens", len(info.RTokens)), zap.Error(err))
	}
	log.Info("user banned", zap.String("login", info.Login), zap.Int("tokens", len(info.RTokens)))
}

func (c *BanChecker) banLoginOnly(ctx context.Context, log *zap.Logger, userId uuid.UUID, ttl time.Duration) {
	info, err := c.auth.GetUserInfo(ctx, userId)
	if err != nil {
		log.Error("failed to get user info for ban", zap.Error(err))
		return
	}
	err = c.repoLogins.PutKey(ctx, info.Login, ttl)
	if err != nil {
		log.Error("failed to ban login", zap.Error(err))
		return
	}
	log.Info("user banned without active tokens", zap.String("login", info.Login))
}
