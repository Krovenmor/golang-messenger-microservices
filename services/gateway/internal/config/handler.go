package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type MainHandlerConfig struct {
	AuthServiceURL   string
	MsgServiceURL    string
	WsServiceURL     string
	WebServiceURL    string
	StatusServiceURL string
}

type MiddlewareConfig struct {
	BanTooManyRequestsDurationUser time.Duration
	BanTooManyRequestsDurationIp   time.Duration

	LimitRateIp  float64
	LimitBurstIp int

	TokenLenght     int
	HeaderLength    int
	PayloadLength   int
	SignatureLength int
}

type InfraConfig struct {
	CacheSize           int
	TickerCleanerTiming time.Duration
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

func GetMiddlewareConfig() (*MiddlewareConfig, error) {
	r := config.NewConfigReader()

	conf := MiddlewareConfig{
		BanTooManyRequestsDurationUser: r.GetDuration("BAN_DUR_TOO_MANY_REQS_USER"),
		BanTooManyRequestsDurationIp:   r.GetDuration("BAN_DUR_TOO_MANY_REQS_IP"),

		LimitRateIp:  r.GetFloat64("LIMIT_RATE_IP"),
		LimitBurstIp: r.GetInt("LIMIT_BURST_IP"),

		TokenLenght:     r.GetInt("TOKEN_LENGTH"),
		HeaderLength:    r.GetInt("TOKEN_HEADER_LENGTH"),
		PayloadLength:   r.GetInt("TOKEN_PAYLOAD_LENGTH"),
		SignatureLength: r.GetInt("TOKEN_SIGNATURE_LENGTH"),
	}

	return &conf, r.Err()
}

func GetInfraConfig() (*InfraConfig, error) {
	r := config.NewConfigReader()

	conf := InfraConfig{
		CacheSize:           r.GetInt("INFRA_CACHE_SIZE"),
		TickerCleanerTiming: r.GetDuration("INFRA_TICKER_TIMING"),
	}

	return &conf, r.Err()
}
