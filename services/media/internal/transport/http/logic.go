package http

import (
	"MyMessenger/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h *Handler) SaveAvatar(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, reqBodyLimit)

	img, _, _, err := GetDecodedImage(r, h.conf.AvatarsSizeLimit, h.conf.AvatarsSizeLimit)
	if err != nil {
		log.Printf("SaveAvatar: trouble with GetDecodedImage, err: %q", err)
		msg, code := errHttp(err)
		http.Error(w, msg, code)
		return
	}

	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("SaveAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	photoId, err := h.serv.SaveAvatar(r.Context(), userId, img)
	if err != nil {
		log.Printf("SaveAvatar: trouble with SaveImage(), err: %q", err)
		msg, code := errHttp(err)
		http.Error(w, msg, code)
		return
	}

	resp := NewPhotoResponseBody{PhotoID: photoId.String()}
	utils.Send(w, &resp)
}

func (h *Handler) DelAvatar(w http.ResponseWriter, r *http.Request) {
	avatarId, err := utils.GetUuidFromPath(r, "uuid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("DelAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = h.serv.DeleteAvatar(r.Context(), userId, avatarId)
	if err != nil {
		msg, code := errHttp(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("DelAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	info, err := h.serv.GetProfileInfo(r.Context(), userId)
	if err != nil {
		msg, code := errHttp(err)
		http.Error(w, msg, code)
		return
	}
	toSend := ToProfileResponseBody(info)
	utils.Send(w, toSend)
}

func (h *Handler) GetSavedFiles(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("DelAvatar: trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	vals := r.URL.Query()
	q, err := strconv.Atoi(vals.Get("q"))
	if err != nil {
		q = -1
	}
	fromId, err := uuid.Parse(vals.Get("from"))
	if err != nil {
		fromId = uuid.Nil
	}

	info, err := h.serv.GetProfileMediaInfo(r.Context(), userId, fromId, q)
	if err != nil {
		msg, code := errHttp(err)
		http.Error(w, msg, code)
		return
	}
	toSend := ToMediaInfoSliceResponseBody(info)
	utils.Send(w, &toSend)
}
