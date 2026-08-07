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
	r := config.NewConfigReader()

	conf := AuthConfig{
		MinPassLength:   r.GetInt("MIN_PASS_LENGTH"),
		MaxPassLength:   r.GetInt("MAX_PASS_LENGTH"),
		MinLoginLength:  r.GetInt("MIN_LOGIN_LENGTH"),
		MaxLoginLength:  r.GetInt("MAX_LOGIN_LENGTH"),
		AccessTokenTTL:  r.GetDuration("JWT_ACCESS_TTL"),
		RefreshTokenTTL: r.GetDuration("JWT_REFRESH_TTL"),
		PrvKey:          r.GetBytes("JWT_PRVKEY_PATH"),
		PubKey:          r.GetBytes("JWT_PUBKEY_PATH"),
	}

	return conf, r.Err()
}
