package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/profile/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

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

func (h *Handler) PostAvatar(w http.ResponseWriter, r *http.Request) {
	avatarId, err := utils.GetUuidFromPath(r, targetKey)
	if err != nil {
		http.Error(w, "Bad Photo ID", http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("PostAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "Internal", http.StatusInternalServerError)
		return
	}
	err = h.msg.AddAvatarToProfile(context.Background(), userId, avatarId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DelAvatar(w http.ResponseWriter, r *http.Request) {
	avatarId, err := utils.GetUuidFromPath(r, targetKey)
	if err != nil {
		http.Error(w, "Bad Photo ID", http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("DelAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "Internal", http.StatusInternalServerError)
		return
	}
	err = h.msg.DelAvatarFromProfile(context.Background(), userId, avatarId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
