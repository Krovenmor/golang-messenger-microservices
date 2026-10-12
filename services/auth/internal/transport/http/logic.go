package http

import (
	"MyMessenger/pkg/httpx"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/domain"
	"net/http"

	"go.uber.org/zap"
)

func sendTokens(w http.ResponseWriter, tokens *domain.Tokens) {
	httpx.WriteJSON(w, http.StatusOK, TokensResponseBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken})
}

func (h *Handler) SendCode(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[SendCodeRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	exp, retry, err := h.auth.SendCodeEmail(r.Context(), body.Email)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, SendCodeResponseBody{ExpiresIn: exp, RetryAfter: retry})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[RegisterRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	err = h.auth.Register(r.Context(), body.Login, body.Password, body.Email, body.Code)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[LoginRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	ctx := r.Context()
	if h.banChecker.CheckLoginBanned(ctx, body.Login) {
		h.log.Ctx(ctx).Info("login attempt with banned login", zap.String("login", body.Login))
		h.writeError(w, r, domain.ErrBanned)
		return
	}
	tokens, err := h.auth.LogIn(ctx, body.Login, body.Password)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[TokensUpdateRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	ctx := r.Context()
	if h.banChecker.CheckTokenBanned(ctx, body.RToken) {
		h.log.Ctx(ctx).Info("tokens update with banned refresh token")
		h.writeError(w, r, domain.ErrBanned)
		return
	}
	tokens, err := h.auth.UpdateTokens(ctx, body.RToken)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUuidFromContext(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	info, err := h.auth.GetUserInfo(r.Context(), userId)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, AccountInfoResponseBody{Login: info.Login, Email: info.Email})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
