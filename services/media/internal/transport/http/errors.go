package http

import (
	"MyMessenger/services/media/internal/service"
	"errors"
	"log"
	"net/http"
)

var (
	ErrBadSize  = errors.New("bad size")
	ErrBadFile  = errors.New("bad file")
	ErrInteranl = errors.New("internal")
)

func errHttp(err error) (string, int) {
	switch {
	case errors.Is(err, ErrBadFile):
		return "bad size", http.StatusBadRequest
	case errors.Is(err, ErrBadFile):
		return "bad file", http.StatusBadRequest
	case errors.Is(err, ErrInteranl):
		return "internal", http.StatusInternalServerError
	case errors.Is(err, service.ErrNotEnoughSpace):
		return "not enough space", http.StatusBadRequest
	case errors.Is(err, service.ErrNotFound):
		return "not found", http.StatusBadRequest
	default:
		log.Printf("Not known error: %q", err)
		return "internal", http.StatusInternalServerError
	}
}
