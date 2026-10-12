package http

import (
	"MyMessenger/pkg/httpx"
	"MyMessenger/services/auth/internal/domain"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var badRequestErrors = []error{
	domain.ErrRegistrationNotFound,
	domain.ErrInvalidCode,
	domain.ErrCodeAlreadySent,
	domain.ErrLoginTaken,
	domain.ErrEmailTaken,
	domain.ErrInvalidCredentials,
	domain.ErrInvalidToken,
	domain.ErrUserNotFound,
}

var errorFields = map[error]httpx.Fields{
	domain.ErrInvalidCode: {"code": {Reason: "invalid"}},
	domain.ErrLoginTaken:  {"login": {Reason: "taken"}},
	domain.ErrEmailTaken:  {"email": {Reason: "taken"}},
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

func fieldsOf(err error) httpx.Fields {
	var decodeErr *httpx.DecodeError
	if errors.As(err, &decodeErr) {
		return decodeErr.Fields()
	}
	for target, fields := range errorFields {
		if errors.Is(err, target) {
			return fields
		}
	}
	return nil
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	status := statusOf(err)
	if status != http.StatusInternalServerError {
		httpx.WriteFieldsError(w, status, err.Error(), fieldsOf(err))
		return
	}
	if !errors.Is(err, domain.ErrInternal) {
		h.log.Ctx(r.Context()).Error("unexpected error", zap.Error(err))
	}
	httpx.WriteError(w, status, domain.ErrInternal.Error())
}
