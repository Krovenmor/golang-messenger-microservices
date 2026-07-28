package http

import (
	"MyMessenger/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

func recv[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var toResv T
	err := utils.Recv(r, &toResv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &toResv, nil
}

func getUuidFromPath(req *http.Request) (uuid.UUID, error) {
	UUID := req.PathValue("uuid")
	if UUID == "" {
		return uuid.Nil, errors.New("Empty UUID")
	}
	cUUID, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.Nil, errors.New("Bad UUID")
	}
	return cUUID, nil
}

func getUUIDQueryParam(vals url.Values, key string) (uuid.UUID, error) {
	val := vals.Get(key)
	if val == "" {
		return uuid.Nil, fmt.Errorf("You must provide %q query param", key)
	}

	valConv, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Not uuid value in %q query param", key)
	}

	return valConv, nil
}

func getIntQueryParam(vals url.Values, key string) (int, error) {
	val := vals.Get(key)
	if val == "" {
		return -1, fmt.Errorf("You must provide %q query param", key)
	}

	valConv, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("Not integer value in %q query param", key)
	}

	return valConv, nil
}

func (h *Handler) NewProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := recv[ProfileBody](w, r)
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

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := h.msg.GetProfile(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if profile == nil {
		http.Error(w, "Trouble with getting profile", http.StatusInternalServerError)
		return
	}
	utils.Send(w, FromServiceProfile(profile))
}

func (h *Handler) NewChat(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := recv[NewChatIncomeBody](w, r)
	if err != nil {
		return
	}
	chatId, err := h.msg.CreateNewChat(r.Context(), userId, body.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Send(w, &NewChatResponseBody{ChatId: chatId})
}

func (h *Handler) PostMessage(w http.ResponseWriter, r *http.Request) {
	chatId, err := getUuidFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := recv[PostMessageIncomeBody](w, r)
	if err != nil {
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mId, err := h.msg.PostMessage(r.Context(), chatId, *ToServiceMsg(body, userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Send(w, &PostMessageResponseBody{MessageId: mId})
}

func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatId, err := getUuidFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vals := r.URL.Query()
	if len(vals) != 2 {
		http.Error(w, "You must provide 'from' and 'q' query params", http.StatusBadRequest)
		return
	}

	from, err := getUUIDQueryParam(vals, "from")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	q, err := getIntQueryParam(vals, "q")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

	utils.Send(w, &msgs)
}
