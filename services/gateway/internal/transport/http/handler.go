package http

import (
	"MyMessenger/services/gateway/internal/config"
	"MyMessenger/services/gateway/internal/middlware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type MainHandler struct {
	middleware *middlware.Middleware
	conf       *config.MainHandlerConfig
}

func NewMainHandler(mw *middlware.Middleware, conf *config.MainHandlerConfig) *MainHandler {
	return &MainHandler{
		middleware: mw,
		conf:       conf,
	}
}

func (h *MainHandler) RegisterRoutes(m *http.ServeMux) (http.Handler, error) {

	var err error
	getParsed := func(toParse string) *httputil.ReverseProxy {
		if err != nil {
			return nil
		}
		var pUrl *url.URL
		pUrl, err = url.Parse(toParse)
		if err != nil {
			return nil
		}
		proxy := httputil.NewSingleHostReverseProxy(pUrl)
		return proxy
	}

	protected := func(pattern string, handler http.Handler) {
		m.Handle(pattern, h.middleware.FullMiddleware(handler))
	}

	authProxy := getParsed(h.conf.AuthServiceURL)
	msgProxy := getParsed(h.conf.MsgServiceURL)
	wsProxy := getParsed(h.conf.WsServiceURL)
	statProxy := getParsed(h.conf.StatusServiceURL)

	if err != nil {
		return nil, err
	}

	protected("/api/msg/", msgProxy)
	protected("/api/status/", statProxy)

	m.Handle("/api/ws", h.middleware.QueryParamMiddleware(wsProxy, "token"))
	m.Handle("/api/auth/", h.middleware.LimitMiddleware(authProxy))

	return m, nil
}
