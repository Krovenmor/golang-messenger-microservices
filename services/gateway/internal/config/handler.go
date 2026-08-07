package config

import "MyMessenger/pkg/config"

type MainHandlerConfig struct {
	AuthServiceURL   string
	MsgServiceURL    string
	WsServiceURL     string
	WebServiceURL    string
	StatusServiceURL string
}

func GetMainHandlerConfig() (*MainHandlerConfig, error) {
	r := config.NewConfigReader()

	conf := MainHandlerConfig{
		AuthServiceURL:   r.GetString("AUTH_URL"),
		MsgServiceURL:    r.GetString("MSG_URL"),
		WsServiceURL:     r.GetString("WS_URL"),
		WebServiceURL:    r.GetString("WEB_URL"),
		StatusServiceURL: r.GetString("STAT_URL"),
	}

	return &conf, r.Err()
}
