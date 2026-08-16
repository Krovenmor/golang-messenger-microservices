package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type AuthConfig struct {
	PrvKey []byte

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	VerificationCodeTTL time.Duration
}

type HashConfig struct {
	Memory      int
	Iterations  int
	Parallelism int
	SaltLength  int
	KeyLength   int
}

type RedisConfig struct {
	RepoLoginPrefix string
	RepoTokenPrefix string

	CacheTTLPrefix string
}

func GetAuthConfig() (*AuthConfig, error) {
	r := config.NewConfigReader()

	conf := AuthConfig{
		AccessTokenTTL:  r.GetDuration("JWT_ACCESS_TTL"),
		RefreshTokenTTL: r.GetDuration("JWT_REFRESH_TTL"),
		PrvKey:          r.GetBytes("JWT_PRVKEY_PATH"),

		VerificationCodeTTL: r.GetDuration("VERIFICATION_CODES_TTL"),
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

func GetRedisConfig() (*RedisConfig, error) {
	r := config.NewConfigReader()

	conf := RedisConfig{
		RepoLoginPrefix: r.GetString("BAN_REPO_LOGIN_PREFIX"),
		RepoTokenPrefix: r.GetString("BAN_REPO_TOKEN_PREFIX"),

		CacheTTLPrefix: r.GetString("CACHE_TTL_PREFIX"),
	}

	return &conf, r.Err()
}
