package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/service"
	"errors"
	"log"
	"net/http"
)

func sendTokens(w http.ResponseWriter, tokens *service.Tokens) {
	tB := TokensResponseBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken}
	utils.Send(w, &tB)
}

func (h *Handler) ToStatusCode(err error) int {
	switch {
	case errors.Is(err, service.ErrBadData):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInternal):
		return http.StatusInternalServerError
	default:
		log.Printf("ToStatusCode: not known error: %q", err)
		return http.StatusBadRequest
	}
}

func (h *Handler) SendCode(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[SendCodeRequestBody](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	exp, retry, err := h.auth.SendCodeEmail(r.Context(), body.Email)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	resp := SendCodeResponseBody{ExpiresIn: exp, RetryAfter: retry}
	utils.Send(w, &resp)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[RegisterRequestBody](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	err = h.auth.Register(r.Context(), body.Login, body.Password, body.Email, body.Code)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[LoginRequestBody](r)
	if err != nil {
		log.Printf("utils.Recv err: %q", err)
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	ctx := r.Context()
	if h.banChecker.CheckLoginBanned(ctx, body.Login) {
		log.Printf("LogIn: banned login: %q", body.Login)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	tokens, err := h.auth.LogIn(ctx, body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[TokensUpdateRequestBody](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	ctx := r.Context()
	if h.banChecker.CheckTokenBanned(ctx, body.RToken) {
		log.Printf("Update: banned token")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	tokens, err := h.auth.UpdateTokens(r.Context(), body.RToken)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		log.Printf("GetAccountInfo: Trouble with GetUuidFromContext, err: %q", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	info, err := h.auth.GetUserInfo(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	resp := AccountInfoResponseBody{Login: info.Login, Email: info.Email}
	utils.Send(w, &resp)
}
