package config

import "MyMessenger/pkg/config"

type ServConfig struct {
	Address string
}

func GetServConfig() (ServConfig, error) {
	var conf ServConfig
	addr, err := config.GetEnvVar("LISTEN_ADDRESS")
	if err != nil {
		return conf, err
	}
	port, err := config.GetEnvVar("LISTEN_PORT")
	if err != nil {
		return conf, err
	}
	return ServConfig{
		Address: addr + ":" + port,
	}, nil
}
