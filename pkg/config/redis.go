package config

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type RedisChannelsConfig struct {
	UserStatusPattern string
	UserEventsPattern string
	UserBanChannel    string
	ChatEventsPattern string
}

func GetRedisConfig() (*RedisConfig, error) {
	r := NewConfigReader()

	conf := RedisConfig{
		Address:  r.GetString("REDIS_ADDRESS"),
		Password: r.GetString("REDIS_PASSWORD"),
		DB:       r.GetInt("REDIS_DB"),
	}

	if r.err != nil {
		return nil, r.err
	}

	return &conf, nil
}

func GetRedisChannelsConfig() (*RedisChannelsConfig, error) {
	r := NewConfigReader()

	conf := RedisChannelsConfig{
		UserStatusPattern: r.GetString("REDIS_USER_STATUS_PATTERN"),
		UserEventsPattern: r.GetString("REDIS_USER_EVENTS_PATTERN"),
		UserBanChannel:    r.GetString("REDIS_USER_BAN_CHANNEL"),
		ChatEventsPattern: r.GetString("REDIS_CHAT_EVENTS_PATTERN"),
	}

	if r.err != nil {
		return nil, r.err
	}

	return &conf, nil
}
