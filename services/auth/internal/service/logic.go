package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/pkg/logger"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/infra/security/code"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JwtAuth struct {
	repo AuthRepo
	conf *config.AuthConfig

	jwtChecker TokenChecker
	gen        TokenGenerator
	hash       AuthHasher

	cache AuthTTLCache
	pub   Publisher

	log *logger.Logger
}

func NewJwtAuth(repo AuthRepo, cache AuthTTLCache, pub Publisher, checker TokenChecker, gen TokenGenerator, hash AuthHasher, conf *config.AuthConfig, log *logger.Logger) *JwtAuth {
	return &JwtAuth{
		repo:       repo,
		conf:       conf,
		jwtChecker: checker,
		gen:        gen,
		hash:       hash,

		cache: cache,
		pub:   pub,

		log: log,
	}
}

func (a *JwtAuth) Register(ctx context.Context, login, password, email, code string) error {
	login, email = domain.NormalizeLogin(login), domain.NormalizeEmail(email)
	log := a.log.Ctx(ctx).With(zap.String("login", login), zap.String("email", email))

	rCode, err := a.cache.Get(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("no verification code for email")
		return domain.ErrInvalidCode
	}
	if err != nil {
		log.Error("failed to get verification code", zap.Error(err))
		return domain.ErrInternal
	}
	if rCode != code {
		log.Debug("verification code mismatch")
		return domain.ErrInvalidCode
	}
	err = a.repo.IsUserExists(ctx, login)
	if err == nil {
		log.Debug("user already exists")
		return domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		log.Error("failed to check user existence", zap.Error(err))
		return domain.ErrInternal
	}
	hashed, err := a.hash.Hash(password)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))
		return domain.ErrInternal
	}
	userId := uuid.New()
	err = a.repo.AddNewUser(ctx, userId, login, hashed, email)
	if errors.Is(err, domain.ErrAlreadyExists) {
		log.Debug("login or email already taken")
		return domain.ErrUserAlreadyExists
	}
	if err != nil {
		log.Error("failed to add user", zap.Error(err))
		return domain.ErrInternal
	}
	err = a.cache.Del(ctx, email)
	if err != nil {
		log.Warn("failed to delete verification code", zap.Error(err))
	}
	log.Info("user registered", zap.Stringer("user_id", userId))
	return nil
}

func (a *JwtAuth) newTokens(ctx context.Context, userId uuid.UUID) (*domain.Tokens, error) {
	log := a.log.Ctx(ctx).With(zap.Stringer("user_id", userId))

	aToken, err := a.gen.GenToken(userId, a.conf.AccessTokenTTL, j.Access)
	if err != nil {
		log.Error("failed to generate access token", zap.Error(err))
		return nil, domain.ErrInternal
	}
	rToken, err := a.gen.GenToken(userId, a.conf.RefreshTokenTTL, j.Refresh)
	if err != nil {
		log.Error("failed to generate refresh token", zap.Error(err))
		return nil, domain.ErrInternal
	}
	err = a.repo.SaveRefresh(ctx, userId, utils.Hash(rToken), time.Now().UTC().Add(a.conf.RefreshTokenTTL))
	if err != nil {
		log.Error("failed to save refresh token", zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &domain.Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

func (a *JwtAuth) LogIn(ctx context.Context, login, password string) (*domain.Tokens, error) {
	login = domain.NormalizeLogin(login)
	log := a.log.Ctx(ctx).With(zap.String("login", login))

	userId, hashed, err := a.repo.GetUser(ctx, login)
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("user not found")
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		log.Error("failed to get user", zap.Error(err))
		return nil, domain.ErrInternal
	}
	err = a.hash.Comp(password, hashed)
	if err != nil {
		log.Debug("password mismatch", zap.Error(err))
		return nil, domain.ErrInvalidCredentials
	}
	err = a.repo.DeleteExpiredRefreshTokens(ctx, userId)
	if err != nil {
		log.Warn("failed to delete expired refresh tokens", zap.Stringer("user_id", userId), zap.Error(err))
	}
	return a.newTokens(ctx, userId)
}

func (a *JwtAuth) UpdateTokens(ctx context.Context, rToken string) (*domain.Tokens, error) {
	log := a.log.Ctx(ctx)

	claims, err := a.jwtChecker.IsValidToken(rToken)
	if err != nil {
		log.Debug("invalid refresh token", zap.Error(err))
		return nil, domain.ErrInvalidToken
	}
	if claims.Ttype != j.Refresh {
		log.Debug("token is not a refresh token", zap.String("type", string(claims.Ttype)))
		return nil, domain.ErrInvalidToken
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Error("failed to parse token subject", zap.Error(err))
		return nil, domain.ErrInternal
	}
	log = log.With(zap.Stringer("user_id", userId))
	err = a.repo.DeleteRefresh(ctx, userId, utils.Hash(rToken))
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("refresh token not found")
		return nil, domain.ErrInvalidToken
	}
	if err != nil {
		log.Error("failed to delete refresh token", zap.Error(err))
		return nil, domain.ErrInternal
	}
	return a.newTokens(ctx, userId)
}

func (a *JwtAuth) GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error) {
	log := a.log.Ctx(ctx).With(zap.Stringer("user_id", userId))

	info, err := a.repo.GetUserInfo(ctx, userId)
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("user not found")
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		log.Error("failed to get user info", zap.Error(err))
		return nil, domain.ErrInternal
	}
	return info, nil
}

func (a *JwtAuth) GetUserTokens(ctx context.Context, userId uuid.UUID) (*domain.UserTokens, error) {
	return a.repo.GetUserTokens(ctx, userId)
}

func (a *JwtAuth) SendCodeEmail(ctx context.Context, email string) (int, int, error) {
	email = domain.NormalizeEmail(email)
	log := a.log.Ctx(ctx).With(zap.String("email", email))

	vCode, err := code.GenerateNewCode()
	if err != nil {
		log.Error("failed to generate verification code", zap.Error(err))
		return -1, -1, domain.ErrInternal
	}
	ttl := a.conf.VerificationCodeTTL
	err = a.cache.Put(ctx, email, vCode, ttl)
	if errors.Is(err, domain.ErrAlreadyExists) {
		log.Debug("verification code already sent")
		return -1, -1, domain.ErrCodeAlreadySent
	}
	if err != nil {
		log.Error("failed to save verification code", zap.Error(err))
		return -1, -1, domain.ErrInternal
	}
	err = a.pub.PublishEmailVerification(ctx, email, vCode)
	if err != nil {
		log.Error("failed to publish email verification", zap.Error(err))
		delErr := a.cache.Del(ctx, email)
		if delErr != nil {
			log.Warn("failed to delete verification code after failed publish", zap.Error(delErr))
		}
		return -1, -1, domain.ErrInternal
	}
	seconds := int(ttl.Seconds())
	return seconds, seconds, nil
}
