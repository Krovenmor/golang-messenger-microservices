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

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type HashConfig struct {
	Memory      int
	Iterations  int
	Parallelism int
	SaltLength  int
	KeyLength   int
}

func GetAuthConfig() (*AuthConfig, error) {
	r := config.NewConfigReader()

	conf := AuthConfig{
		MinPassLength:   r.GetInt("MIN_PASS_LENGTH"),
		MaxPassLength:   r.GetInt("MAX_PASS_LENGTH"),
		MinLoginLength:  r.GetInt("MIN_LOGIN_LENGTH"),
		MaxLoginLength:  r.GetInt("MAX_LOGIN_LENGTH"),
		AccessTokenTTL:  r.GetDuration("JWT_ACCESS_TTL"),
		RefreshTokenTTL: r.GetDuration("JWT_REFRESH_TTL"),
		PrvKey:          r.GetBytes("JWT_PRVKEY_PATH"),
	}

	return &conf, r.Err()
}

func GetHashConfig() (*HashConfig, error) {
	r := config.NewConfigReader()

	conf := HashConfig{
		Memory:      r.GetInt("HASH_MEMORY"),
		Iterations:  r.GetInt("HASH_ITERATIONS"),
		Parallelism: r.GetInt("HASH_PARALLELISM"),
		SaltLength:  r.GetInt("HASH_SALT_LENGTH"),
		KeyLength:   r.GetInt("HASH_KEY_LENGTH"),
	}

	return &conf, r.Err()
}
