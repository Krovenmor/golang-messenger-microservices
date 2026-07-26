package service

import (
	"MyMessenger/services/auth/internal/config"
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthServiceImpl struct {
	repo AuthRepo
	conf config.AuthConfig

	PubKey *rsa.PublicKey
	PrvKey *rsa.PrivateKey
}

func NewAuth(repo AuthRepo, conf config.AuthConfig) *AuthServiceImpl {

	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(conf.PrvKey)
	if err != nil {
		log.Fatalf("NewAuth(): Trouble with parsing prvKey: %q", err.Error())
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(conf.PubKey)
	if err != nil {
		log.Fatalf("NewAuth(): Trouble with parsing pubKey: %q", err.Error())
	}

	return &AuthServiceImpl{repo: repo, conf: conf, PubKey: pubKey, PrvKey: prvKey}
}

func (a *AuthServiceImpl) checkLoginPass(login, password string) error {
	lLen := len(login)
	if lLen < a.conf.MinLoginLength {
		return fmt.Errorf("Login is too short, min len: %d", a.conf.MinLoginLength)
	}
	if lLen > a.conf.MaxLoginLength {
		return fmt.Errorf("Login is too big, max len: %d", a.conf.MaxLoginLength)
	}
	pLen := len(password)
	if pLen < a.conf.MinPassLength {
		return fmt.Errorf("Password is too short, min len: %d", a.conf.MinPassLength)
	}
	if pLen > a.conf.MaxPassLength {
		return fmt.Errorf("Password is too big, max len: %d", a.conf.MaxPassLength)
	}
	return nil
}

func (a *AuthServiceImpl) Register(ctx context.Context, login, password string) error {
	login, password = strings.TrimSpace(login), strings.TrimSpace(password)
	err := a.checkLoginPass(login, password)
	if err != nil {
		return err
	}
	userId := uuid.New()
	err = a.repo.AddNewUser(ctx, userId, login, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthServiceImpl) newTokens(ctx context.Context, userId uuid.UUID) (*Tokens, error) {
	aToken, err := GenToken(userId, a.conf.AccessTokenTTL, a.PrvKey)
	if err != nil {
		return nil, err
	}
	rToken, err := GenToken(userId, a.conf.RefreshTokenTTL, a.PrvKey)
	if err != nil {
		return nil, err
	}
	err = a.repo.SaveRefresh(ctx, userId, rToken, time.Now().UTC().Add(a.conf.RefreshTokenTTL))
	if err != nil {
		return nil, err
	}
	return &Tokens{AccessToken: aToken, RefreshToken: rToken}, nil
}

func (a *AuthServiceImpl) LogIn(ctx context.Context, login, password string) (*Tokens, error) {
	login, password = strings.TrimSpace(login), strings.TrimSpace(password)
	err := a.checkLoginPass(login, password)
	if err != nil {
		return nil, err
	}
	userId, err := a.repo.GetUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return a.newTokens(ctx, userId)
}

func (a *AuthServiceImpl) IsValidAccess(ctx context.Context, aToken string) error {
	claims, err := IsValidToken(aToken, a.PubKey)
	if err != nil {
		return err
	}
	subj, err := claims.GetSubject()
	if err != nil {
		return err
	}
	userId, err := uuid.Parse(subj)
	if err != nil {
		return err
	}
	if userId == uuid.Nil {
		return fmt.Errorf("Nil UUID")
	}
	return nil
}

func (a *AuthServiceImpl) UpdateTokens(ctx context.Context, rToken string) (*Tokens, error) {
	claims, err := IsValidToken(rToken, a.PubKey)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}
	err = a.repo.FindRefresh(ctx, userId, rToken)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("rToken expired at %v", claims.ExpiresAt.Time)
	}
	nTokens, err := a.newTokens(ctx, userId)
	if err != nil {
		return nil, err
	}
	_ = a.repo.DeleteRefresh(ctx, userId, rToken)
	return nTokens, nil
}
