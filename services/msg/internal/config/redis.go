package config

import "MyMessenger/pkg/config"

type RedisPatternConfig struct {
	PubUserPattern string
	PubChatPattern string
}

func GetRedisPatternConfig() (*RedisPatternConfig, error) {
	pubUser, err := config.GetEnvVar("REDIS_PUB_USER_PATTERN")
	if err != nil {
		return nil, err
	}
	pubChat, err := config.GetEnvVar("REDIS_PUB_CHAT_PATTERN")
	if err != nil {
		return nil, err
	}
	return &RedisPatternConfig{
		PubUserPattern: pubUser,
		PubChatPattern: pubChat,
	}, nil
}
