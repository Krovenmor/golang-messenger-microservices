package http

import (
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/msg/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var (
	chatIdKey    = "chatid"
	messageIdKey = "messageid"
	targetKey    = "target"
	emojiKey     = "emoji"
)

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

	profilePath := " /api/msg/profile"
	targetPath := fmt.Sprintf("/{%s}", targetKey)

	protected("POST"+profilePath, h.NewProfile)
	protected("GET"+profilePath, h.GetProfilePrivate)
	protected("GET"+profilePath+targetPath, h.GetProfilePublic)
	protected("GET"+profilePath+"/chats", h.GetProfileChats)
	protected("GET"+profilePath+"/chats/full", h.GetProfileChatsExtended)

	protected("POST"+profilePath+"/avatar"+targetPath, h.PostAvatar)
	protected("DELETE"+profilePath+"/avatar"+targetPath, h.DelAvatar)

	chatPath := " /api/msg/chat"
	chatWithIdPath := chatPath + fmt.Sprintf("/{%s}", chatIdKey)

	protected("POST"+chatPath, h.NewChat)
	protected("POST"+chatWithIdPath, h.PostMessage)
	protected("GET"+chatWithIdPath, h.GetChatMessages)
	protected("GET"+chatWithIdPath+"/info", h.GetChatInfo)

	chatAndMessagePath := chatWithIdPath + fmt.Sprintf("/message/{%s}", messageIdKey)

	protected("GET"+chatAndMessagePath, h.GetMessage)
	protected("PUT"+chatAndMessagePath, h.ChangeMessage)
	protected("DELETE"+chatAndMessagePath, h.DeleteMessage)

	protected("GET /api/msg/reactions", h.GetSupportedReactions)

	reactionPath := chatAndMessagePath + "/reaction"
	reactionEmojiPath := reactionPath + fmt.Sprintf("/{%s}", emojiKey)

	protected("POST"+reactionPath, h.PostReaction)
	protected("DELETE"+reactionEmojiPath, h.DeleteReaction)

	return m
}
