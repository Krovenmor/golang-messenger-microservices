package httpx

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type FieldError struct {
	Reason string `json:"reason"`
	Param  string `json:"param,omitempty"`
}

type Fields map[string]FieldError

type ErrorResponse struct {
	Error  string `json:"error"`
	Fields Fields `json:"fields,omitempty"`
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

func WriteFieldsError(w http.ResponseWriter, status int, msg string, fields Fields) {
	WriteJSON(w, status, ErrorResponse{Error: msg, Fields: fields})
}
