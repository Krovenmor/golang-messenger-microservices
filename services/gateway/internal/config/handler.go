package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type MainHandlerConfig struct {
	AuthServiceURL   string
	MsgServiceURL    string
	WsServiceURL     string
	StatusServiceURL string
	MediaServiceURL  string
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
	CacheSize int
}

func GetMainHandlerConfig() (*MainHandlerConfig, error) {
	r := config.NewConfigReader()

	conf := MainHandlerConfig{
		AuthServiceURL:   r.GetString("AUTH_URL"),
		MsgServiceURL:    r.GetString("MSG_URL"),
		WsServiceURL:     r.GetString("WS_URL"),
		StatusServiceURL: r.GetString("STAT_URL"),
		MediaServiceURL:  r.GetString("MEDIA_URL"),
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

		TokenLenght:     r.GetInt("JWT_ACCESS_TOKEN_LEN"),
		HeaderLength:    r.GetInt("TOKEN_HEADER_LENGTH"),
		PayloadLength:   r.GetInt("TOKEN_PAYLOAD_LENGTH"),
		SignatureLength: r.GetInt("TOKEN_SIGNATURE_LENGTH"),
	}

	return &conf, r.Err()
}

func GetInfraConfig() (*InfraConfig, error) {
	r := config.NewConfigReader()

	conf := InfraConfig{
		CacheSize: r.GetInt("CACHE_SIZE"),
	}

	return &conf, r.Err()
}
