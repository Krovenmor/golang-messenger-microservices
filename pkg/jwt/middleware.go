package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Authenticator struct {
	checker JWTChecker
}

func NewAuthenticator(checker *JWTChecker) *Authenticator {
	return &Authenticator{
		checker: *checker,
	}
}

type contextKey string

const UserIdKey contextKey = "userUUID"

func (a *Authenticator) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)

		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			log.Printf("No bearer auth was included")
			http.Error(w, "Bearer auth required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			log.Printf("Bearer auth without token")
			http.Error(w, "Empty token", http.StatusUnauthorized)
			return
		}

		claims, err := a.checker.IsValidAccess(token)
		if err != nil {
			log.Printf("False token")
			http.Error(w, fmt.Sprintf("Trouble with token, err: %q", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, claims.Subject)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
