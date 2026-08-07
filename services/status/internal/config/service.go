package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type ServiceConfig struct {
	EntriesTTL time.Duration
}

type SubInfoConfig struct {
	SubPattern string
}

func GetServiceConfig() (*ServiceConfig, error) {
	ttl, err := config.GetEnvVarDuration("REDIS_STATUSES_TTL")
	if err != nil {
		return nil, err
	}
	return &ServiceConfig{
		EntriesTTL: ttl,
	}, nil
}

func GetSubInfoConfig() (*SubInfoConfig, error) {
	prefix, err := config.GetEnvVar("REDIS_USER_STATUS_PATTERN")
	if err != nil {
		return nil, err
	}
	return &SubInfoConfig{
		SubPattern: prefix,
	}, nil
}
