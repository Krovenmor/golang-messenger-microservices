package http

import (
	"MyMessenger/pkg/httpx"
	"MyMessenger/services/auth/internal/domain"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var badRequestErrors = []error{
	domain.ErrInvalidCode,
	domain.ErrCodeAlreadySent,
	domain.ErrUserAlreadyExists,
	domain.ErrInvalidCredentials,
	domain.ErrInvalidToken,
	domain.ErrUserNotFound,
}

func statusOf(err error) int {
	var decodeErr *httpx.DecodeError
	if errors.As(err, &decodeErr) {
		return http.StatusBadRequest
	}
	if errors.Is(err, domain.ErrBanned) {
		return http.StatusForbidden
	}
	for _, target := range badRequestErrors {
		if errors.Is(err, target) {
			return http.StatusBadRequest
		}
	}
	return http.StatusInternalServerError
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := statusOf(err)
	if status != http.StatusInternalServerError {
		httpx.WriteError(w, status, err.Error())
		return
	}
	if !errors.Is(err, domain.ErrInternal) {
		h.log.Ctx(r.Context()).Error("unexpected error", zap.Error(err))
	}
	httpx.WriteError(w, status, domain.ErrInternal.Error())
}
