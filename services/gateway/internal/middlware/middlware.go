package middlware

import (
	stdconfig "MyMessenger/pkg/config"
	"MyMessenger/pkg/jwt"
	"MyMessenger/services/gateway/internal/config"
	"fmt"
	"log"
	"net"
	"net/http"

	"go.uber.org/fx"
	"golang.org/x/time/rate"
)

type Middleware struct {
	repoBansId   MiddlewareRepo
	repoBansAddr MiddlewareRepo
	repoTokens   MiddlewareRepo

	addrCache MiddlewareCache

	checker *jwt.JWTChecker

	conf *config.MiddlewareConfig
}

func NewMiddleware(
	lf fx.Lifecycle,
	repoFab func(lf fx.Lifecycle) MiddlewareRepo,
	cacheFab func() MiddlewareCache,
	sub Subscriber,
	checker *jwt.JWTChecker,
	confC *stdconfig.RedisChannelsConfig, confM *config.MiddlewareConfig,
) (*Middleware, error) {

	repoBansId, repoBansAddr, repoTokens := repoFab(lf), repoFab(lf), repoFab(lf)
	if repoBansAddr == nil || repoBansId == nil || repoTokens == nil {
		return nil, fmt.Errorf("trouble with NewRepo, null repo")
	}

	addrC := cacheFab()
	if addrC == nil {
		return nil, fmt.Errorf("trouble with NewCache, null cache")
	}

	m := &Middleware{
		repoBansId:   repoBansId,
		repoBansAddr: repoBansAddr,
		repoTokens:   repoTokens,

		addrCache: addrC,

		conf: confM,

		checker: checker,
	}

	registerReader(lf, sub, confC.UserBanChannel, repoBansId, m.getBanTime)
	return m, nil
}

func getClientIP(r *http.Request) string {
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		log.Printf("X-Real-IP: %q", xri)
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Err with SplitHostPort, err: %q", err)
		return r.RemoteAddr
	}
	log.Printf("SplitHostPort: %q", host)
	return host
}

func (m *Middleware) LimitMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("LimitMiddleware: Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)

		addr := getClientIP(r)
		if m.repoBansAddr.IsExists(addr) {
			log.Printf("Banned addr: %q", addr)
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		limiter, isExists := m.addrCache.Get(addr)
		if !isExists {
			log.Printf("NewLimiter for addr: %q", addr)
			limiter = rate.NewLimiter(rate.Limit(m.conf.LimitRateIp), m.conf.LimitBurstIp)
			m.addrCache.Put(addr, limiter)
		}

		if !limiter.Allow() {
			log.Printf("New ban for addr: %q", addr)
			m.repoBansAddr.Put(addr, m.getBanTime(TooManyRequestsIp))
			http.Error(w, "", http.StatusTooManyRequests)
			return
		}

		log.Printf("LimitMiddleware passed for addr: %q", addr)
		h.ServeHTTP(w, r)
	})
}

func (m *Middleware) QueryParamMiddleware(h http.Handler, param string) http.Handler {
	return m.LimitMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("QueryParamMiddleware: Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)

			token := r.URL.Query().Get(param)
			err := m.checkToken(token)
			if err != nil {
				log.Printf("QueryParamMiddleware: failed for token: %q", token)
				http.Error(w, err.Error(), getStatusFromError(err))
				return
			}

			log.Printf("QueryParamMiddleware: passed for token: %q", token)
			h.ServeHTTP(w, r)
		}),
	)
}

func (m *Middleware) FullMiddleware(h http.Handler) http.Handler {
	return m.LimitMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("FullMiddleware: Income request to: %q, from: %q", r.URL.Path, r.RemoteAddr)

			token, err := jwt.GetBearerToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			err = m.checkToken(token)
			if err != nil {
				log.Printf("FullMiddleware: failed for token: %q", token)
				http.Error(w, err.Error(), getStatusFromError(err))
				return
			}

			log.Printf("FullMiddleware: passed for token: %q", token)
			h.ServeHTTP(w, r)
		}),
	)
}
