package config

import "MyMessenger/pkg/config"

type WsConfig struct {
	WsReadLimit int64
}

type RedisPatternConfig struct {
	PatternPub string
	PatternSub string
}

func GetRedisPatternConfig() (*RedisPatternConfig, error) {
	pub, err := config.GetEnvVar("REDIS_PUB_PATTERN")
	if err != nil {
		return nil, err
	}
	sub, err := config.GetEnvVar("REDIS_SUB_PATTERN")
	if err != nil {
		return nil, err
	}
	return &RedisPatternConfig{
		PatternPub: pub,
		PatternSub: sub,
	}, nil
}

func GetWsConfig() (*WsConfig, error) {
	limit, err := config.GetEnvVarInt("WS_READ_LIMIT")
	if err != nil {
		return nil, err
	}
	return &WsConfig{
		WsReadLimit: int64(limit),
	}, nil
}
