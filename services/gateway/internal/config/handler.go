package config

import "MyMessenger/pkg/config"

type MainHandlerConfig struct {
	AuthServiceURL string
	MsgServiceURL  string
	WsServiceURL   string
}

func GetMainHandlerConfig() (*MainHandlerConfig, error) {
	var (
		conf MainHandlerConfig
		err  error
	)

	get := func(envKey string) string {
		if err != nil {
			return ""
		}
		var val string
		val, err = config.GetEnvVar(envKey)
		return val
	}

	conf = MainHandlerConfig{
		AuthServiceURL: get("AUTH_URL"),
		MsgServiceURL:  get("MSG_URL"),
		WsServiceURL:   get("WS_URL"),
	}

	if err != nil {
		return nil, err
	}

	return &conf, nil
}
