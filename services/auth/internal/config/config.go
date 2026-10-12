package config

import (
	"time"
)

type AuthConfig struct {
	PrvKey []byte `env:"JWT_PRVKEY_PATH,file"`

	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TTL"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TTL"`

	VerificationCodeTTL time.Duration `env:"VERIFICATION_CODES_TTL"`
}

type HashConfig struct {
	Memory      int `env:"HASH_MEMORY"`
	Iterations  int `env:"HASH_ITERATIONS"`
	Parallelism int `env:"HASH_PARALLELISM"`
	SaltLength  int `env:"HASH_SALT_LENGTH"`
	KeyLength   int `env:"HASH_KEY_LENGTH"`
}

type RedisConfig struct {
	RepoLoginPrefix string `env:"BAN_REPO_LOGIN_PREFIX"`
	RepoTokenPrefix string `env:"BAN_REPO_TOKEN_PREFIX"`

	CacheTTLPrefix string `env:"CACHE_TTL_PREFIX"`
}

type ValidationConfig struct {
	PassMinLength int `env:"PASS_MIN_LENGTH"`
	PassMaxLength int `env:"PASS_MAX_LENGTH"`

	LoginMinLength int `env:"LOGIN_MIN_LENGTH"`
	LoginMaxLength int `env:"LOGIN_MAX_LENGTH"`

	RefreshTokenLen int `env:"REFRESH_TOKEN_LEN"`
}
