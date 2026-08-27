package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/media/internal/config"
	"MyMessenger/services/media/internal/service"
	"net/http"
)

const (
	reqBodyLimit = 5 * 1024 * 1024
)

type Handler struct {
	auth *jwt.Authenticator
	serv service.MediaService

	conf *config.HandlerConfig
}

func NewHandler(auth *jwt.Authenticator, serv service.MediaService, conf *config.HandlerConfig) *Handler {
	return &Handler{auth: auth, serv: serv, conf: conf}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	protected := func(pattern string, handlerFunc http.HandlerFunc) {
		m.Handle(pattern, h.auth.Middleware(handlerFunc))
	}

	protected("GET /api/media/profile", h.GetProfile)
	protected("GET /api/media/profile/files", h.GetSavedFiles)
	protected("POST /api/media/public/avatar", h.SaveAvatar)
	protected("DELETE /api/media/public/avatar/{uuid}", h.DelAvatar)

	return m
}
