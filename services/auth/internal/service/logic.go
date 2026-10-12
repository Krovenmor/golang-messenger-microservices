package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/infra/security/code"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type JwtAuth struct {
	repo AuthRepo
	conf *config.AuthConfig

	jwtChecker TokenChecker
	gen        TokenGenerator
	hash       AuthHasher

	cache AuthTTLCache
	pub   Publisher
}

func NewJwtAuth(repo AuthRepo, cache AuthTTLCache, pub Publisher, checker TokenChecker, gen TokenGenerator, hash AuthHasher, conf *config.AuthConfig) *JwtAuth {
	return &JwtAuth{
		repo:       repo,
		conf:       conf,
		jwtChecker: checker,
		gen:        gen,
		hash:       hash,

		cache: cache,
		pub:   pub,
	}
}

func (a *JwtAuth) Register(ctx context.Context, login, password, email, code string) error {
	login, email = domain.NormalizeLogin(login), domain.NormalizeEmail(email)
	rCode, err := a.cache.Get(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		log.Printf("Register: no verification code for %q", email)
		return domain.ErrBadData
	}
	if err != nil {
		log.Printf("Register: Trouble with cache.Get, err: %q", err)
		return domain.ErrInternal
	}
	if rCode != code {
		log.Printf("Register: Codes not equal: %v != %v", rCode, code)
		return domain.ErrBadData
	}
	err = a.repo.IsUserExists(ctx, login)
	if err == nil {
		log.Printf("Register: user already exists: %q", login)
		return domain.ErrBadData
	}
	if !errors.Is(err, domain.ErrNotFound) {
		log.Printf("Register: Trouble with IsUserExists, login: %q, err: %q", login, err)
		return domain.ErrInternal
	}
	hashed, err := a.hash.Hash(password)
	if err != nil {
		log.Printf("Register: Trouble with hashing password, err: %q", err)
		return domain.ErrInternal
	}
	userId := uuid.New()
	err = a.repo.AddNewUser(ctx, userId, login, hashed, email)
	if errors.Is(err, domain.ErrAlreadyExists) {
		log.Printf("Register: login or email already taken, login: %q, email: %q", login, email)
		return domain.ErrBadData
	}
	if err != nil {
		log.Printf("Register: Trouble with AddNewUser, login: %q, err: %q", login, err)
		return domain.ErrInternal
	}
	err = a.cache.Del(ctx, email)
	if err != nil {
		log.Printf("Register: trouble with deleting email, login:%q, email:%q, err:%q", login, email, err)
	}
	return nil
}

func (a *JwtAuth) newTokens(ctx context.Context, userId uuid.UUID) (*domain.Tokens, error) {
	aToken, err := a.gen.GenToken(userId, a.conf.AccessTokenTTL, j.Access)
	if err != nil {
		log.Printf("newTokens: Trouble with generating aToken for user %q, err: %q", userId.String(), err)
		return nil, domain.ErrInternal
	}
	rToken, err := a.gen.GenToken(userId, a.conf.RefreshTokenTTL, j.Refresh)
	if err != nil {
		log.Printf("newTokens: Trouble with generating rToken for user %q, err: %q", userId.String(), err)
		return nil, domain.ErrInternal
	}
	err = a.repo.SaveRefresh(ctx, userId, utils.Hash(rToken), time.Now().UTC().Add(a.conf.RefreshTokenTTL))
	if err != nil {
		log.Printf("newTokens: Trouble with saving rToken for user %q, err: %q", userId.String(), err)
		return nil, domain.ErrInternal
	}
	return &domain.Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

func (a *JwtAuth) LogIn(ctx context.Context, login, password string) (*domain.Tokens, error) {
	login = domain.NormalizeLogin(login)
	userId, hashed, err := a.repo.GetUser(ctx, login)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrBadData
	}
	if err != nil {
		log.Printf("LogIn: Trouble with GetUser, login: %q, err: %q", login, err)
		return nil, domain.ErrInternal
	}
	err = a.hash.Comp(password, hashed)
	if err != nil {
		log.Printf("LogIn: Trouble with hash.Comp(), login: %q, err: %q", login, err)
		return nil, domain.ErrBadData
	}
	err = a.repo.DeleteExpiredRefreshTokens(ctx, userId)
	if err != nil {
		log.Printf("LogIn: Trouble with DeleteExpiredRefreshTokens, userId: %q, err: %q", userId.String(), err)
	}
	return a.newTokens(ctx, userId)
}

func (a *JwtAuth) UpdateTokens(ctx context.Context, rToken string) (*domain.Tokens, error) {
	claims, err := a.jwtChecker.IsValidToken(rToken)
	if err != nil {
		log.Printf("UpdateTokens: Not valid refresh token, err: %q", err)
		return nil, domain.ErrBadData
	}
	if claims.Ttype != j.Refresh {
		log.Printf("UpdateTokens: Not a refresh token, type: %q", claims.Ttype)
		return nil, domain.ErrBadData
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with uuid.Parse, err: %q", err)
		return nil, domain.ErrInternal
	}
	err = a.repo.DeleteRefresh(ctx, userId, utils.Hash(rToken))
	if errors.Is(err, domain.ErrNotFound) {
		log.Printf("UpdateTokens: refresh token not found, userId: %q", userId.String())
		return nil, domain.ErrBadData
	}
	if err != nil {
		log.Printf("UpdateTokens: Trouble with DeleteRefresh, userId: %q, err: %q", userId.String(), err)
		return nil, domain.ErrInternal
	}
	return a.newTokens(ctx, userId)
}

func (a *JwtAuth) GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error) {
	info, err := a.repo.GetUserInfo(ctx, userId)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrBadData
	}
	if err != nil {
		log.Printf("GetUserInfo: Trouble with repo.GetUserInfo, userId: %q, err: %q", userId.String(), err)
		return nil, domain.ErrInternal
	}
	return info, nil
}

func (a *JwtAuth) GetUserTokens(ctx context.Context, userId uuid.UUID) (*domain.UserTokens, error) {
	return a.repo.GetUserTokens(ctx, userId)
}

func (a *JwtAuth) SendCodeEmail(ctx context.Context, email string) (int, int, error) {
	email = domain.NormalizeEmail(email)
	vCode, err := code.GenerateNewCode()
	if err != nil {
		log.Printf("SendCodeEmail: Trouble with GenerateNewCode, err: %q", err)
		return -1, -1, domain.ErrInternal
	}
	ttl := a.conf.VerificationCodeTTL
	err = a.cache.Put(ctx, email, vCode, ttl)
	if errors.Is(err, domain.ErrAlreadyExists) {
		log.Printf("SendCodeEmail: Already send code to %q", email)
		return -1, -1, domain.ErrBadData
	}
	if err != nil {
		log.Printf("SendCodeEmail: Trouble with Put, err: %q", err)
		return -1, -1, domain.ErrInternal
	}
	err = a.pub.PublishEmailVerification(ctx, email, vCode)
	if err != nil {
		log.Printf("SendCodeEmail: Trouble with PublishEmailVerification, err: %q", err)
		delErr := a.cache.Del(ctx, email)
		if delErr != nil {
			log.Printf("SendCodeEmail: Trouble with Del after failed publish, err: %q", delErr)
		}
		return -1, -1, domain.ErrInternal
	}
	seconds := int(ttl.Seconds())
	return seconds, seconds, nil
}
