package config

type ServConfig struct {
	Address string
}

func GetServConfig() (ServConfig, error) {
	var conf ServConfig
	addr, err := GetEnvVar("LISTEN_ADDRESS")
	if err != nil {
		return conf, err
	}
	port, err := GetEnvVar("LISTEN_PORT")
	if err != nil {
		return conf, err
	}
	return ServConfig{
		Address: addr + ":" + port,
	}, nil
}
