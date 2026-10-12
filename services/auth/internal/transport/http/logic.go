package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/domain"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func sendTokens(w http.ResponseWriter, tokens *domain.Tokens) {
	tB := TokensResponseBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken}
	utils.Send(w, &tB)
}

func (h *Handler) ToStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrBadData):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrInternal):
		return http.StatusInternalServerError
	default:
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
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	ctx := r.Context()
	if h.banChecker.CheckLoginBanned(ctx, body.Login) {
		h.log.Ctx(ctx).Info("login attempt with banned login", zap.String("login", body.Login))
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
		h.log.Ctx(ctx).Info("tokens update with banned refresh token")
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
		h.log.Ctx(r.Context()).Error("failed to get user id from context", zap.Error(err))
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

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
