package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra/security/code"
	"context"
	"errors"
	"log"
	"strings"
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
	login, email = strings.ToLower(login), strings.ToLower(email)
	rCode, err := a.cache.Get(ctx, email)
	if err != nil {
		log.Printf("Register: Get, err: %q", err)
		return ErrBadData
	}
	if rCode != code {
		log.Printf("Register: Codes not equal: %v != %v", rCode, code)
		return ErrBadData
	}
	err = a.repo.IsUserExists(ctx, login)
	if err == nil {
		log.Printf("Register: user already exists: %q", login)
		return ErrBadData
	}
	hashed, err := a.hash.Hash(password)
	if err != nil {
		log.Printf("Register: Trouble with hashing password, err: %q", err)
		return ErrInternal
	}
	userId := uuid.New()
	err = a.repo.AddNewUser(ctx, userId, login, hashed, email)
	if err != nil {
		log.Printf("Register: trouble with adding new user after repo.IsExists, login:%q, err:%q", login, err)
		return ErrBadData
	}
	err = a.cache.Del(ctx, email)
	if err != nil {
		log.Printf("Register: trouble with deleting email, login:%q, email:%q, err:%q", login, email, err)
	}
	return nil
}

func (a *JwtAuth) newTokens(ctx context.Context, userId uuid.UUID) (*Tokens, error) {
	aToken, err := a.gen.GenToken(userId, a.conf.AccessTokenTTL, j.Access)
	if err != nil {
		log.Printf("newTokens: Trouble with generating aToken for user %q, err: %q", userId.String(), err)
		return nil, ErrInternal
	}
	rToken, err := a.gen.GenToken(userId, a.conf.RefreshTokenTTL, j.Refresh)
	if err != nil {
		log.Printf("newTokens: Trouble with generating rToken for user %q, err: %q", userId.String(), err)
		return nil, ErrInternal
	}
	err = a.repo.SaveRefresh(ctx, userId, utils.Hash(rToken), time.Now().UTC().Add(a.conf.RefreshTokenTTL))
	if err != nil {
		log.Printf("newTokens: Trouble with saving rToken for user %q, err: %q", userId.String(), err)
		return nil, ErrInternal
	}
	return &Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

func (a *JwtAuth) LogIn(ctx context.Context, login, password string) (*Tokens, error) {
	login = strings.ToLower(login)
	userId, hashed, err := a.repo.GetUser(ctx, login)
	if err != nil {
		return nil, ErrBadData
	}
	err = a.hash.Comp(password, hashed)
	if err != nil {
		log.Printf("LogIn: Trouble with hash.Comp(), login: %q, err: %q", login, err)
		return nil, ErrBadData
	}
	a.repo.DeleteExpiredRefreshTokens(ctx, userId)
	return a.newTokens(ctx, userId)
}

func (a *JwtAuth) UpdateTokens(ctx context.Context, rToken string) (*Tokens, error) {
	claims, err := a.jwtChecker.IsValidToken(rToken)
	if err != nil {
		log.Printf("UpdateTokens: Not valid refresh token, err: %q", err)
		return nil, ErrBadData
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with uuid.Parse, err: %q", err)
		return nil, ErrInternal
	}
	hashedRToken := utils.Hash(rToken)
	err = a.repo.FindRefresh(ctx, userId, hashedRToken)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with FindRefresh, userId: %q, err: %q", userId.String(), err)
		return nil, ErrBadData
	}
	if time.Now().UTC().After(claims.ExpiresAt.Time) {
		log.Printf("UpdateTokens: Expired refresh token, err: %q", err)
		return nil, ErrBadData
	}
	nTokens, err := a.newTokens(ctx, userId)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with newTokens, err: %q", err)
		return nil, err
	}
	err = a.repo.DeleteRefresh(ctx, userId, hashedRToken)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with DeleteRefresh, userId: %q, err: %q", userId.String(), err)
	}
	return nTokens, nil
}

func (a *JwtAuth) GetUserInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error) {
	return a.repo.GetUserInfo(ctx, userId)
}

func (a *JwtAuth) GetUserTokens(ctx context.Context, userId uuid.UUID) (*UserTokens, error) {
	return a.repo.GetUserTokens(ctx, userId)
}

func (a *JwtAuth) SendCodeEmail(ctx context.Context, email string) (int, int, error) {
	email = strings.ToLower(email)
	_, err := a.cache.Get(ctx, email)
	if err == nil {
		log.Printf("SendCodeEmail: Already send code to %q", email)
		return -1, -1, ErrBadData
	}
	vCode, err := code.GenerateNewCode()
	if err != nil {
		log.Printf("SendCodeEmail: Trouble with GenerateNewCode, err: %q", err)
		return -1, -1, ErrInternal
	}
	err = a.pub.PublishEmailVerification(ctx, email, vCode)
	if err != nil {
		log.Printf("SendCodeEmail: Trouble with PublishEmailVerification, err: %q", err)
		return -1, -1, ErrInternal
	}
	ttl := a.conf.VerificationCodeTTL
	err = a.cache.Put(ctx, email, vCode, ttl)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return -1, -1, ErrBadData
		}
		log.Printf("SendCodeEmail: Trouble with Put, err: %q", err)
		return -1, -1, ErrInternal
	}
	seconds := int(ttl.Seconds())
	return seconds, seconds, nil
}
