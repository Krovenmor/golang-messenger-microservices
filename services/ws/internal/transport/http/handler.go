package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/ws/internal/transport/ws"
	"log"
	"net/http"
)

type MainHandler struct {
	checker   *jwt.JWTChecker
	wsHandler *ws.WSHandler
}

func NewMainHandler(checker *jwt.JWTChecker, wsHandler *ws.WSHandler) *MainHandler {
	return &MainHandler{
		checker:   checker,
		wsHandler: wsHandler,
	}
}

func (h *MainHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	if len(vals) != 1 {
		http.Error(w, "You must provide valid AccessToken in query param", http.StatusUnauthorized)
		return
	}
	accessToken := vals.Get("token")
	if accessToken == "" {
		http.Error(w, "You must provide valid AccessToken in query param", http.StatusUnauthorized)
		return
	}
	claims, err := h.checker.IsValidAccess(accessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId, err := h.checker.GetUserIdFromClaims(claims)
	if err != nil {
		log.Printf("Error with GetUserIdFromClaims, err: %q", err.Error())
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.wsHandler.HandleConnection(w, r, userId)
}

func (h *MainHandler) RegisterRoutes(m *http.ServeMux) (http.Handler, error) {
	m.HandleFunc("/api/ws/", h.HandleWS)
	return m, nil
}
