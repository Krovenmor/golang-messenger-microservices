package config

import (
	"MyMessenger/pkg/config"
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

func GetConfig() (*AuthConfig, error) {
	minPl, err := getParsedInt("MIN_PASS_LENGTH")
	if err != nil {
		return nil, err
	}
	maxPl, err := getParsedInt("MAX_PASS_LENGTH")
	if err != nil {
		return nil, err
	}
	minLl, err := getParsedInt("MIN_LOGIN_LENGTH")
	if err != nil {
		return nil, err
	}
	maxLl, err := getParsedInt("MAX_LOGIN_LENGTH")
	if err != nil {
		return nil, err
	}

	return &AuthConfig{
		MinPassLength:  minPl,
		MinLoginLength: minLl,
		MaxPassLength:  maxPl,
		MaxLoginLength: maxLl,
	}, nil
}
