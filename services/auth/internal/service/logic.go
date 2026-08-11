package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/config"
	"context"
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
}

func NewJwtAuth(repo AuthRepo, checker TokenChecker, gen TokenGenerator, hash AuthHasher, conf *config.AuthConfig) *JwtAuth {
	return &JwtAuth{
		repo:       repo,
		conf:       conf,
		jwtChecker: checker,
		gen:        gen,
		hash:       hash,
	}
}

func (a *JwtAuth) Register(ctx context.Context, login, password string) error {
	err := a.repo.IsUserExists(ctx, login)
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
	err = a.repo.AddNewUser(ctx, userId, login, hashed)
	if err != nil {
		log.Printf("Register: trouble with adding new user after repo.IsExists, login:%q, err:%q", login, err)
		return ErrBadData
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
		return nil, ErrBadData
	}
	nTokens, err := a.newTokens(ctx, userId)
	if err != nil {
		return nil, err
	}
	err = a.repo.DeleteRefresh(ctx, userId, hashedRToken)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with DeleteRefresh, userId: %q, err: %q", userId.String(), err)
	}
	return nTokens, nil
}

func (a *JwtAuth) GetInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error) {
	return a.repo.GetUserInfo(ctx, userId)
}
