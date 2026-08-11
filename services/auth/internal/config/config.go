package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type AuthConfig struct {
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

type BanConfig struct {
	RepoLoginPrefix string
	RepoTokenPrefix string
}

func GetAuthConfig() (*AuthConfig, error) {
	r := config.NewConfigReader()

	conf := AuthConfig{
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

func GetBanConfig() (*BanConfig, error) {
	r := config.NewConfigReader()

	conf := BanConfig{
		RepoLoginPrefix: r.GetString("BAN_REPO_LOGIN_PREFIX"),
		RepoTokenPrefix: r.GetString("BAN_REPO_TOKEN_PREFIX"),
	}

	return &conf, r.Err()
}
