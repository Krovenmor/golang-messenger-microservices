package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/msg/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var chatIdKey = "chatid"
var messageIdKey = "messageid"
var targetKey = "target"

type Handler struct {
	msg  service.MessageService
	auth *jwt.Authenticator
}

func NewHandler(msgService service.MessageService, auth *jwt.Authenticator) *Handler {
	auth.SetExternalCheckFunc(func(w http.ResponseWriter, r *http.Request, c *jwt.TokenClaims) error {
		chatUUID := r.PathValue(chatIdKey)
		if chatUUID == "" {
			return nil
		}
		chatId, err := uuid.Parse(chatUUID)
		if err != nil {
			err := fmt.Errorf("invalid chatId UUID: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		userId, err := uuid.Parse(c.Subject)
		if err != nil {
			log.Printf("ExtCheckerFunc: Trouble with getting uuid from TokenClaims, err: %q", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return err
		}
		err = msgService.IsProfileInChat(r.Context(), userId, chatId)
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return err
		}
		return nil
	})
	return &Handler{msg: msgService, auth: auth}
}

func (h *Handler) RegisterRoutes(m *http.ServeMux) http.Handler {
	protected := func(pattern string, handlerFunc http.HandlerFunc) {
		m.Handle(pattern, h.auth.Middleware(handlerFunc))
	}

	protected("POST /api/msg/profile/new", h.NewProfile)
	protected("GET /api/msg/profile", h.GetProfilePrivate)
	protected(fmt.Sprintf("GET /api/msg/profile/{%s}", targetKey), h.GetProfilePublic)
	protected("GET /api/msg/profile/chats", h.GetProfileChats)
	protected("GET /api/msg/profile/chats/full", h.GetProfileChatsExtended)

	protected("POST /api/msg/chat/new", h.NewChat)
	protected(fmt.Sprintf("POST /api/msg/chat/{%s}", chatIdKey), h.PostMessage)
	protected(fmt.Sprintf("GET /api/msg/chat/{%s}", chatIdKey), h.GetChatMessages)
	protected(fmt.Sprintf("GET /api/msg/chat/{%s}/info", chatIdKey), h.GetChatInfo)

	protected(fmt.Sprintf("GET /api/msg/chat/{%s}/message/{%s}", chatIdKey, messageIdKey), h.GetMessage)
	protected(fmt.Sprintf("PUT /api/msg/chat/{%s}/message/{%s}", chatIdKey, messageIdKey), h.ChangeMessage)
	protected(fmt.Sprintf("DELETE /api/msg/chat/{%s}/message/{%s}", chatIdKey, messageIdKey), h.DeleteMessage)

	return m
}
