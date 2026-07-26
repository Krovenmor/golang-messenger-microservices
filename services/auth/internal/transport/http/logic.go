package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/auth/internal/service"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body AuthIncomeBody
	err := utils.Recv(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.auth.Register(r.Context(), body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func sendTokens(w http.ResponseWriter, tokens *service.Tokens) {
	tB := TokensOutcomeBody{AToken: tokens.AccessToken, RToken: tokens.RefreshToken}
	utils.Send(w, &tB)
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	var body AuthIncomeBody
	err := utils.Recv(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := h.auth.LogIn(r.Context(), body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if tokens == nil {
		http.Error(w, "Trouble with recieving tokens", http.StatusInternalServerError)
		return
	}
	sendTokens(w, tokens)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var body TokensUpdateIncome
	err := utils.Recv(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := h.auth.UpdateTokens(r.Context(), body.RToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if tokens == nil {
		http.Error(w, "Trouble with updating tokens", http.StatusInternalServerError)
		return
	}
	sendTokens(w, tokens)
}
