package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/service"
	"errors"
	"net/http"
)

func (h *Handler) ToStatusCode(err error) int {
	switch {
	case errors.Is(err, service.ErrBadData):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInternal):
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[AuthIncomeBody](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	err = h.auth.Register(r.Context(), body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func sendTokens(w http.ResponseWriter, tokens *service.Tokens) {
	tB := TokensOutcomeBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken}
	utils.Send(w, &tB)
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[AuthIncomeBody](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	tokens, err := h.auth.LogIn(r.Context(), body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := utils.Recv[TokensUpdateIncome](r)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	tokens, err := h.auth.UpdateTokens(r.Context(), body.RToken)
	if err != nil {
		http.Error(w, err.Error(), h.ToStatusCode(err))
		return
	}
	sendTokens(w, tokens)
}
