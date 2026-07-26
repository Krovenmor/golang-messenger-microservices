package config

import (
	"MyMessenger/pkg/config"
	"os"
	"strconv"
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

func getParsedInt(key string) (int, error) {
	val, err := config.GetEnvVar(key)
	if err != nil {
		return -1, err
	}
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return -1, err
	}
	return iVal, nil
}

func getParsedDuration(key string) (time.Duration, error) {
	val, err := config.GetEnvVar(key)
	if err != nil {
		return -1, err
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		return -1, err
	}
	return dur, nil
}

func getKey(key string) ([]byte, error) {
	val, err := config.GetEnvVar(key)
	if err != nil {
		return []byte{}, err
	}
	keyAr, err := os.ReadFile(val)
	if err != nil {
		return []byte{}, err
	}
	return keyAr, nil
}

func GetAuthConfig() (AuthConfig, error) {
	var conf AuthConfig
	minPl, err := getParsedInt("MIN_PASS_LENGTH")
	if err != nil {
		return conf, err
	}
	maxPl, err := getParsedInt("MAX_PASS_LENGTH")
	if err != nil {
		return conf, err
	}
	minLl, err := getParsedInt("MIN_LOGIN_LENGTH")
	if err != nil {
		return conf, err
	}
	maxLl, err := getParsedInt("MAX_LOGIN_LENGTH")
	if err != nil {
		return conf, err
	}
	aTTL, err := getParsedDuration("JWT_ACCESS_TTL")
	if err != nil {
		return conf, err
	}
	rTTL, err := getParsedDuration("JWT_REFRESH_TTL")
	if err != nil {
		return conf, err
	}
	prvKey, err := getKey("JWT_PRVKEY_PATH")
	if err != nil {
		return conf, err
	}
	pubKey, err := getKey("JWT_PUBKEY_PATH")
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
