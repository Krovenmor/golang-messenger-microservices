package ban

import (
	"MyMessenger/pkg/broker"
	stdconf "MyMessenger/pkg/config"
	"MyMessenger/pkg/redis/reader"
	"MyMessenger/pkg/utils"
	"strings"

	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/service"
	"context"
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

// {"type": "ban_ev", "payload": {"userId": "98d199b7-0513-4418-8acd-ccb3d4d67575", "ttl": 100000000000}}

func NewBanChecker(lf fx.Lifecycle, repoFab func(prefix string) MiddlewareRepo, sub Subscriber, auth service.AuthService,
	banConf *config.BanConfig, redisConf *stdconf.RedisChannelsConfig) *BanChecker {
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
	return c.repoLogins.IsExists(ctx, strings.ToLower(login))
}

func (c *BanChecker) CheckTokenBanned(ctx context.Context, rToken string) bool {
	return c.repoTokens.IsExists(ctx, utils.Hash(rToken))
}

func (c *BanChecker) banUser(ctx context.Context, userId uuid.UUID, ttl time.Duration) {
	info, err := c.auth.GetInfo(ctx, userId)
	if err != nil {
		log.Printf("banUser: trouble with GetInfo(%q), err: %q", userId, err)
		return
	}
	err = c.repoLogins.Put(ctx, info.Login, ttl)
	if err != nil {
		log.Printf("banUser: Trouble with repoLogins.Put(ctx, %v, %d)", userId, ttl)
	}
	err = c.repoTokens.PutKeys(ctx, info.RTokens, ttl)
	if err != nil {
		log.Printf("banUser: Trouble with repoLogins.PutKeys(ctx, %v, %d)", info.RTokens, ttl)
	}
}
