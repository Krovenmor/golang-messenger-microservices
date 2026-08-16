package http

import (
	"MyMessenger/services/auth/internal/service"
	"context"
	"log"
	"net/http"
)

type BanChecker interface {
	CheckLoginBanned(ctx context.Context, login string) bool
	CheckTokenBanned(ctx context.Context, rToken string) bool
}

type Authenticator interface {
	Middleware(h http.Handler) http.Handler
}

type Handler struct {
	auth       service.AuthService
	banChecker BanChecker

	a Authenticator
}

func NewHandler(auth service.AuthService, banChecker BanChecker, a Authenticator) *Handler {
	return &Handler{auth: auth, banChecker: banChecker, a: a}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.HandleFunc("POST /api/auth/send-code", h.SendCode)
	m.HandleFunc("POST /api/auth/register", h.Register)
	m.HandleFunc("POST /api/auth/login", h.LogIn)
	m.HandleFunc("POST /api/auth/update", h.Update)

	m.Handle("GET /api/auth/account", h.a.Middleware(http.HandlerFunc(h.GetAccountInfo)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Income request: %q", r.URL.Path)
		m.ServeHTTP(w, r)
	})
}
