package config

import "MyMessenger/pkg/config"

type RepoConfig struct {
	ConnString string
}

func GetRepoConfig() (RepoConfig, error) {
	var conf RepoConfig
	connString, err := config.GetEnvVar("DB_CONN_STRING")
	if err != nil {
		return conf, err
	}
	return RepoConfig{
		ConnString: connString,
	}, nil
}
