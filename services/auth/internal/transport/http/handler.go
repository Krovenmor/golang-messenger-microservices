package http

import (
	"MyMessenger/pkg/di"
	"MyMessenger/pkg/logger"
	"MyMessenger/services/auth/internal/service"
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
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

	validate *validator.Validate
	log      *logger.Logger
}

func NewHandler(auth service.AuthService, banChecker BanChecker, a Authenticator, validate *validator.Validate, log *logger.Logger) *Handler {
	return &Handler{auth: auth, banChecker: banChecker, a: a, validate: validate, log: log}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.HandleFunc("POST /api/auth/register", h.StartRegistration)
	m.HandleFunc("POST /api/auth/register/confirm", h.ConfirmRegistration)
	m.HandleFunc("POST /api/auth/register/resend", h.ResendCode)
	m.HandleFunc("POST /api/auth/login", h.LogIn)
	m.HandleFunc("POST /api/auth/update", h.Update)
	m.HandleFunc("GET "+di.HealthPath, h.Health)

	m.Handle("GET /api/auth/account", h.a.Middleware(http.HandlerFunc(h.GetAccountInfo)))

	return m
}
