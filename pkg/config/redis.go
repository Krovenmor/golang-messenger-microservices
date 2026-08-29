package config

type RedisConfig struct {
	Address  string
	Password string
	DB       int

	StreamsMaxLen int64
}

type RedisChannelsConfig struct {
	UserStatusPattern string
	UserEventsPattern string

	UserBanRequestChannel   string
	UserBanEventChannel     string
	UserVerificationChannel string

	ChatEventsPattern string

	ProfileStream string
}

func GetRedisConfig() (*RedisConfig, error) {
	r := NewConfigReader()

	conf := RedisConfig{
		Address:  r.GetString("REDIS_ADDRESS"),
		Password: r.GetString("REDIS_PASSWORD"),
		DB:       r.GetInt("REDIS_DB"),

		StreamsMaxLen: int64(r.GetInt("REDIS_STREAMS_MAX_LEN")),
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

		UserBanRequestChannel:   r.GetString("REDIS_USER_BAN_REQUEST_CHANNEL"),
		UserBanEventChannel:     r.GetString("REDIS_USER_BAN_EVENT_CHANNEL"),
		UserVerificationChannel: r.GetString("REDIS_USER_VERIFICATION_CHANNEL"),

		ChatEventsPattern: r.GetString("REDIS_CHAT_EVENTS_PATTERN"),

		ProfileStream: r.GetString("REDIS_PROFILE_STREAM"),
	}

	if r.err != nil {
		return nil, r.err
	}

	return &conf, nil
}
