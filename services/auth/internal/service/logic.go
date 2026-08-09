package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/services/auth/internal/config"
	"context"
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

func (a *JwtAuth) checkLoginPass(login, password string) error {
	lLen := len(login)
	if lLen < a.conf.MinLoginLength {
		log.Printf("checkLoginPass: login is too short, min len: %d", a.conf.MinLoginLength)
		return ErrBadData
	}
	if lLen > a.conf.MaxLoginLength {
		log.Printf("checkLoginPass: login is too big, max len: %d", a.conf.MaxLoginLength)
		return ErrBadData
	}
	pLen := len(password)
	if pLen < a.conf.MinPassLength {
		log.Printf("checkLoginPass: password is too short, min len: %d", a.conf.MinPassLength)
		return ErrBadData
	}
	if pLen > a.conf.MaxPassLength {
		log.Printf("checkLoginPass: password is too big, max len: %d", a.conf.MaxPassLength)
		return ErrBadData
	}
	return nil
}

func (a *JwtAuth) getCheckedLoginPass(login, password string) (string, string, error) {
	login, password = strings.TrimSpace(login), strings.TrimSpace(password)
	login = strings.ToLower(login)
	return login, password, a.checkLoginPass(login, password)
}

func (a *JwtAuth) Register(ctx context.Context, login, password string) error {
	login, password, err := a.getCheckedLoginPass(login, password)
	if err != nil {
		return err
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
	err = a.repo.SaveRefresh(ctx, userId, rToken, time.Now().UTC().Add(a.conf.RefreshTokenTTL))
	if err != nil {
		log.Printf("newTokens: Trouble with saving rToken for user %q, err: %q", userId.String(), err)
		return nil, ErrInternal
	}
	return &Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

func (a *JwtAuth) LogIn(ctx context.Context, login, password string) (*Tokens, error) {
	login, password, err := a.getCheckedLoginPass(login, password)
	if err != nil {
		return nil, err
	}
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
	err = a.repo.FindRefresh(ctx, userId, rToken)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with FindRefresh, userId: %q, err: %q", userId.String(), err)
		return nil, err
	}
	if time.Now().UTC().After(claims.ExpiresAt.Time) {
		return nil, ErrBadData
	}
	nTokens, err := a.newTokens(ctx, userId)
	if err != nil {
		return nil, err
	}
	err = a.repo.DeleteRefresh(ctx, userId, rToken)
	if err != nil {
		log.Printf("UpdateTokens: Trouble with DeleteRefresh, userId: %q, err: %q", userId.String(), err)
	}
	return nTokens, nil
}
