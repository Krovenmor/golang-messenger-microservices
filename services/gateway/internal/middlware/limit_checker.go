package middlware

import (
	"log"
	"net/http"
)

type LimitChecker struct{}

func NewLimitChecker() *LimitChecker {
	return &LimitChecker{}
}

func (l *LimitChecker) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("LimitChecker: Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}
