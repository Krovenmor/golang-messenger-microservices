package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/msg/internal/service"
	"net/http"
)

type Handler struct {
	msg  service.MessageService
	auth *jwt.Authenticator
}

func NewHandler(msgService service.MessageService, auth *jwt.Authenticator) *Handler {
	return &Handler{msg: msgService, auth: auth}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.HandleFunc("POST /api/msg/profile/new", h.NewProfile)
	m.HandleFunc("GET /api/msg/profile", h.GetProfile)

	m.HandleFunc("POST /api/msg/chat/new", h.NewChat)
	m.HandleFunc("POST /api/msg/chat/{uuid}", h.PostMessage)
	m.HandleFunc("GET /api/msg/chat/{uuid}", h.GetChat)

	return h.auth.Middleware(m)
}
