package config

import "MyMessenger/pkg/config"

type WsConfig struct {
	WsReadLimit int64
}

type RedisPatternConfig struct {
	PubPattern string

	SubUserPattern string
	SubChatPattern string
}

type MsgClientConfig struct {
	FullURL string
}

func GetRedisPatternConfig() (*RedisPatternConfig, error) {
	pub, err := config.GetEnvVar("REDIS_PUB_PATTERN")
	if err != nil {
		return nil, err
	}
	subUser, err := config.GetEnvVar("REDIS_SUB_USER_PATTERN")
	if err != nil {
		return nil, err
	}
	subChat, err := config.GetEnvVar("REDIS_SUB_CHAT_PATTERN")
	if err != nil {
		return nil, err
	}
	return &RedisPatternConfig{
		PubPattern:     pub,
		SubUserPattern: subUser,
		SubChatPattern: subChat,
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

func GetMsgClientConfig() (*MsgClientConfig, error) {
	msgBase, err := config.GetEnvVar("MSG_URL")
	if err != nil {
		return nil, err
	}
	msgChatsPath, err := config.GetEnvVar("MSG_GET_CHATS_PATH")
	if err != nil {
		return nil, err
	}
	return &MsgClientConfig{
		FullURL: msgBase + msgChatsPath,
	}, nil
}
