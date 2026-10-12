package http

import (
	"MyMessenger/pkg/httpx"
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/domain"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func parseRegistrationId(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, domain.ErrRegistrationNotFound
	}
	return id, nil
}

func sendTokens(w http.ResponseWriter, tokens *domain.Tokens) {
	httpx.WriteJSON(w, http.StatusOK, TokensResponseBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken})
}

func (h *Handler) StartRegistration(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[StartRegistrationRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	registrationId, ttl, err := h.auth.StartRegistration(r.Context(), body.Login, body.Password, body.Email)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	seconds := int(ttl.Seconds())
	httpx.WriteJSON(w, http.StatusAccepted, StartRegistrationResponseBody{
		RegistrationId: registrationId.String(),
		ExpiresIn:      seconds,
		RetryAfter:     seconds,
	})
}

func (h *Handler) ResendCode(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[ResendCodeRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	registrationId, err := parseRegistrationId(body.RegistrationId)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	ttl, err := h.auth.ResendCode(r.Context(), registrationId)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	seconds := int(ttl.Seconds())
	httpx.WriteJSON(w, http.StatusOK, SendCodeResponseBody{ExpiresIn: seconds, RetryAfter: seconds})
}

func (h *Handler) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.DecodeJSON[ConfirmRegistrationRequestBody](r, h.validate)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	registrationId, err := parseRegistrationId(body.RegistrationId)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	err = h.auth.ConfirmRegistration(r.Context(), registrationId, body.Code)
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
