package config

import "MyMessenger/pkg/config"

type RedisPatternConfig struct {
	PubPattern string
}

func GetRedisPatternConfig() (*RedisPatternConfig, error) {
	pub, err := config.GetEnvVar("REDIS_PUB_PATTERN")
	if err != nil {
		return nil, err
	}
	return &RedisPatternConfig{
		PubPattern: pub,
	}, nil
}
