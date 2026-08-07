package config

import (
	"MyMessenger/pkg/config"
	"time"
)

type ServiceConfig struct {
	EntriesTTL time.Duration
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
