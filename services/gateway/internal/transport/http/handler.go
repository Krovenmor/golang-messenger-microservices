package http

import (
	"MyMessenger/services/gateway/internal/config"
	"MyMessenger/services/gateway/internal/middlware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type MainHandler struct {
	limitChecker *middlware.LimitChecker

	conf *config.MainHandlerConfig
}

func NewMainHandler(limitChecker *middlware.LimitChecker, conf *config.MainHandlerConfig) *MainHandler {

	return &MainHandler{
		limitChecker: limitChecker,
		conf:         conf,
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
		m.Handle(pattern, h.limitChecker.Middleware(handler))
	}

	authProxy := getParsed(h.conf.AuthServiceURL)
	msgProxy := getParsed(h.conf.MsgServiceURL)
	wsProxy := getParsed(h.conf.WsServiceURL)
	webProxy := getParsed(h.conf.WebServiceURL)
	statProxy := getParsed(h.conf.StatusServiceURL)

	if err != nil {
		return nil, err
	}

	protected("/api/auth/", authProxy)
	protected("/api/msg/", msgProxy)
	protected("/api/ws", wsProxy)
	protected("/api/status/", statProxy)

	m.Handle("/", webProxy)

	return m, nil
}
