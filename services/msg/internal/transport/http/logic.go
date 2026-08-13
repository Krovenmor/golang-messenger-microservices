package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/msg/internal/service"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const DEF_QUANTITY_QUERIES = 10

func (h *Handler) NewProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := recv[NewProfileRequestBody](w, r)
	if err != nil {
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.msg.NewProfile(r.Context(), *ToServiceProfile(profile, userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetProfilePrivate(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := h.msg.GetProfileById(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if profile == nil {
		http.Error(w, "Trouble with getting profile", http.StatusInternalServerError)
		return
	}
	utils.Send(w, ToPrivateProfileBody(profile))
}

func (h *Handler) GetProfilePublic(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue(targetKey)
	if userName == "" {
		http.Error(w, fmt.Sprintf("You must provide %q", targetKey), http.StatusBadRequest)
		return
	}
	var profile *service.Profile
	userId, err := uuid.Parse(userName)
	if err != nil {
		profile, err = h.msg.GetProfileByUserName(r.Context(), userName)
	} else {
		profile, err = h.msg.GetProfileById(r.Context(), userId)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if profile == nil {
		http.Error(w, "Trouble with getting profile", http.StatusInternalServerError)
		return
	}
	utils.Send(w, ToPublicProfileBody(profile))
}

func (h *Handler) NewChat(w http.ResponseWriter, r *http.Request) {
	body, err := recv[NewChatRequestBody](w, r)
	if err != nil {
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chatId, err := h.msg.CreateNewChat(r.Context(), userId, body.UserId)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			utils.SendWithStatus(w, &NewChatResponseBody{ChatId: chatId}, http.StatusOK)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendWithStatus(w, &NewChatResponseBody{ChatId: chatId}, http.StatusCreated)
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	chatId, err := getChatUuidFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := recv[PostMessageRequestBody](w, r)
	if err != nil {
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mId, err := h.msg.PostMessage(r.Context(), chatId, *ToServicePostMsg(body, userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Send(w, &PostMessageResponseBody{MessageId: mId})
}

func (h *Handler) GetProfileChats(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chats, err := h.msg.GetChats(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Send(w, &chats)
}

func (h *Handler) GetProfileChatsExtended(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chats, err := h.msg.GetChatsExtended(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chatsToSend := FromServiceFullChatsInfo(chats)
	utils.Send(w, &chatsToSend)
}

func (h *Handler) GetChatInfo(w http.ResponseWriter, r *http.Request) {
	chatId, err := getChatUuidFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chatInfo, err := h.msg.GetChatInfo(r.Context(), chatId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Send(w, FromServiceChatInfo(chatInfo))
}

func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	chatId, err := getChatUuidFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vals := r.URL.Query()
	from, err := getUUIDQueryParam(vals, "from")
	if err != nil {
		from = uuid.Max
	}

	q, err := getIntQueryParam(vals, "q")
	if err != nil {
		q = DEF_QUANTITY_QUERIES
	}

	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msgs, err := h.msg.GetChatHistory(r.Context(), chatId, userId, from, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msgsToSend := FromServiceMessages(msgs)
	utils.Send(w, &msgsToSend)
}

func (h *Handler) GetMessage(w http.ResponseWriter, r *http.Request) {
	chatId, msgId, err := getChatIdAndMsgIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := h.msg.GetMessage(r.Context(), chatId, msgId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if msg == nil {
		log.Printf("msg.GetMessage() returned nil msg")
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	msgToSend := FromServiceMessage(*msg)
	utils.Send(w, &msgToSend)
}

func (h *Handler) ChangeMessage(w http.ResponseWriter, r *http.Request) {
	chatId, msgId, err := getChatIdAndMsgIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := recv[PostMessageRequestBody](w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.msg.RedactMessage(r.Context(), chatId, msgId, *ToServicePostMsg(msg, userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	chatId, msgId, err := getChatIdAndMsgIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.msg.DelMessage(r.Context(), chatId, msgId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetSupportedReactions(w http.ResponseWriter, r *http.Request) {
	reactions, err := h.msg.GetSupportedReactions(r.Context())
	if err != nil {
		log.Printf("Trouble with GetSupportedReactions, err: %q", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	utils.Send(w, &reactions)
}

func (h *Handler) PostReaction(w http.ResponseWriter, r *http.Request) {
	chatId, msgId, err := getChatIdAndMsgIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req, err := utils.Recv[PostReactionRequestBody](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.msg.PostReaction(r.Context(), userId, chatId, msgId, req.Emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteReaction(w http.ResponseWriter, r *http.Request) {
	chatId, msgId, err := getChatIdAndMsgIdFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	emoji, err := getEmojiFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.msg.DelReaction(r.Context(), userId, chatId, msgId, emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
