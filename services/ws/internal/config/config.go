package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type WsConfig struct {
	ReadLimit int64

	LimitRate       float64
	LimitBurst      int
	LimitViolations int

	TickerTiming    time.Duration
	ContextWaitTime time.Duration
}

type MsgClientConfig struct {
	FullURL string
}

func GetWsConfig() (*WsConfig, error) {
	r := config.NewConfigReader()

	conf := WsConfig{
		ReadLimit: int64(r.GetInt("CONN_READ_LIMIT")),

		LimitRate:       r.GetFloat64("WS_LIMIT_RATE"),
		LimitBurst:      r.GetInt("WS_LIMIT_BURST"),
		LimitViolations: r.GetInt("WS_LIMIT_VIOLATIONS"),

		TickerTiming:    r.GetDuration("WORKER_TICKER_TIMING"),
		ContextWaitTime: r.GetDuration("WORKER_CTX_WAIT_TIME"),
	}

	return &conf, r.Err()
}

func GetMsgClientConfig() (*MsgClientConfig, error) {
	r := config.NewConfigReader()

	conf := MsgClientConfig{
		FullURL: r.GetString("MSG_URL") + r.GetString("MSG_GET_CHATS_PATH"),
	}

	return &conf, r.Err()
}
