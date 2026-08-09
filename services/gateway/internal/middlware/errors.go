package middlware

import (
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	ErrInternal = errors.New("internal Error")
	ErrBanned   = errors.New("banned")
	ErrBadToken = errors.New("bad token")

	DefBanTime = time.Minute * 30
)

func getStatusFromError(err error) int {
	switch err {
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrBanned:
		return http.StatusTooManyRequests
	case ErrBadToken:
		return http.StatusUnauthorized

	default:
		log.Printf("Unknown error: %v", err)
		return http.StatusBadRequest
	}
}
