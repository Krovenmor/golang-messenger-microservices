package jwt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	ErrNoAuth  = errors.New("bearer auth required")
	ErrNoToken = errors.New("empty token")
)

type ExtCheckFunc func(w http.ResponseWriter, r *http.Request, c *TokenClaims) error

type Authenticator struct {
	checker  JWTChecker
	extCheck ExtCheckFunc
}

func NewAuthenticator(checker *JWTChecker) *Authenticator {
	return &Authenticator{
		checker: *checker,
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

		log.Printf("Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)

		token, err := GetBearerToken(r)
		if err != nil {
			log.Printf("GetBearerToken: %q", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := a.checker.IsValidAccess(token)
		if err != nil {
			log.Printf("False token")
			http.Error(w, fmt.Sprintf("Trouble with token, err: %q", err), http.StatusUnauthorized)
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
