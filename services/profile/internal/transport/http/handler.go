package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/profile/internal/service"
	"fmt"
	"net/http"
)

var (
	profilePath = " /api/profile"

	targetKey  = "target"
	targetPath = fmt.Sprintf("/{%s}", targetKey)
)

type Handler struct {
	msg  service.ProfileService
	auth *jwt.Authenticator
}

func NewHandler(msgService service.ProfileService, auth *jwt.Authenticator) *Handler {
	return &Handler{msg: msgService, auth: auth}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	protected := func(pattern string, handlerFunc http.HandlerFunc) {
		m.Handle(pattern, h.auth.Middleware(handlerFunc))
	}

	protected("POST"+profilePath, h.NewProfile)
	protected("GET"+profilePath, h.GetProfilePrivate)
	protected("GET"+profilePath+targetPath, h.GetProfilePublic)
	protected("POST"+profilePath+"/batch", h.GetProfilesPublicBatch)

	protected("POST"+profilePath+"/avatar"+targetPath, h.PostAvatar)
	protected("DELETE"+profilePath+"/avatar"+targetPath, h.DelAvatar)

	return m
}
