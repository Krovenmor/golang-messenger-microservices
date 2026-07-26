package http

import (
	"MyMessenger/services/auth/internal/service"
	"log"
	"net/http"
)

type Handler struct {
	auth service.AuthService
}

func NewHandler(auth service.AuthService) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.HandleFunc("POST /api/auth/register", h.Register)
	m.HandleFunc("POST /api/auth/login", h.LogIn)
	m.HandleFunc("POST /api/auth/update", h.Update)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Income request: %q", r.URL.Path)
		m.ServeHTTP(w, r)
	})
}
