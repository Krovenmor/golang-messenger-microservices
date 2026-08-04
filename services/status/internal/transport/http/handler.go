package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/status/internal/service"
	"net/http"
)

type Handler struct {
	auth      *jwt.Authenticator
	stService service.StatusService
}

func NewHandler(auth *jwt.Authenticator, stService service.StatusService) *Handler {
	return &Handler{auth: auth, stService: stService}
}

func (h *Handler) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromPath(r, "uuid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userStatus := h.stService.GetStatus(r.Context(), userId.String())
	utils.Send(w, &userStatus)
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	m.Handle("GET /api/status/{uuid}", h.auth.Middleware(http.HandlerFunc(h.HandleGetStatus)))
	return m
}
