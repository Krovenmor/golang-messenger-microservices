package httpx

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		zap.L().Error("failed to encode response", zap.Error(err))
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}
