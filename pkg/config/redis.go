package config

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func GetRedisConfig() (*RedisConfig, error) {
	addr, err := GetEnvVar("REDIS_ADDRESS")
	if err != nil {
		return nil, err
	}
	pass, err := GetEnvVar("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}
	db, err := GetEnvVarInt("REDIS_DB")
	if err != nil {
		return nil, err
	}
	return &RedisConfig{
		Address:  addr,
		Password: pass,
		DB:       db,
	}, nil
}
