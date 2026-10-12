package jwt

import (
	"MyMessenger/pkg/logger"
	"context"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrNoAuth  = errors.New("bearer auth required")
	ErrNoToken = errors.New("empty token")
)

type ExtCheckFunc func(w http.ResponseWriter, r *http.Request, c *TokenClaims) error

type Authenticator struct {
	checker  JWTChecker
	extCheck ExtCheckFunc
	log      *logger.Logger
}

func NewAuthenticator(checker *JWTChecker, log *logger.Logger) *Authenticator {
	return &Authenticator{
		checker: *checker,
		log:     log,
	}
}

func (a *Authenticator) SetExternalCheckFunc(extChecker ExtCheckFunc) {
	a.extCheck = extChecker
}

type contextKey string

const UserIdKey contextKey = "userUUID"

func GetBearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return "", ErrNoAuth
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" {
		return "", ErrNoToken
	}

	return token, nil
}

func (a *Authenticator) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetBearerToken(r)
		if err != nil {
			a.log.Ctx(r.Context()).Debug("no bearer token", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := a.checker.IsValidAccess(token)
		if err != nil {
			a.log.Ctx(r.Context()).Debug("invalid access token", zap.Error(err))
			http.Error(w, "Bad token", http.StatusUnauthorized)
			return
		}

		if a.extCheck != nil {
			err = a.extCheck(w, r, claims)
			if err != nil {
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIdKey, claims.Subject)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) GetChecker() *JWTChecker {
	return &a.checker
}
