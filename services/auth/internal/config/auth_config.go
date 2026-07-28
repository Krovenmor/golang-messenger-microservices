package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type AuthConfig struct {
	MinPassLength int
	MaxPassLength int

	MinLoginLength int
	MaxLoginLength int

	PrvKey []byte
	PubKey []byte

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func GetAuthConfig() (AuthConfig, error) {
	var conf AuthConfig
	minPl, err := config.GetEnvVarInt("MIN_PASS_LENGTH")
	if err != nil {
		return conf, err
	}
	maxPl, err := config.GetEnvVarInt("MAX_PASS_LENGTH")
	if err != nil {
		return conf, err
	}
	minLl, err := config.GetEnvVarInt("MIN_LOGIN_LENGTH")
	if err != nil {
		return conf, err
	}
	maxLl, err := config.GetEnvVarInt("MAX_LOGIN_LENGTH")
	if err != nil {
		return conf, err
	}
	aTTL, err := config.GetEnvVarDuration("JWT_ACCESS_TTL")
	if err != nil {
		return conf, err
	}
	rTTL, err := config.GetEnvVarDuration("JWT_REFRESH_TTL")
	if err != nil {
		return conf, err
	}
	prvKey, err := config.GetEnvVarBytes("JWT_PRVKEY_PATH")
	if err != nil {
		return conf, err
	}
	pubKey, err := config.GetEnvVarBytes("JWT_PUBKEY_PATH")
	if err != nil {
		return conf, err
	}

	return AuthConfig{
		MinPassLength:   minPl,
		MinLoginLength:  minLl,
		MaxPassLength:   maxPl,
		MaxLoginLength:  maxLl,
		AccessTokenTTL:  aTTL,
		RefreshTokenTTL: rTTL,
		PubKey:          pubKey,
		PrvKey:          prvKey,
	}, nil
}
