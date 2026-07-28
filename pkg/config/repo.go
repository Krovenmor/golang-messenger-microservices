package config

type RepoConfig struct {
	ConnString string
}

func GetRepoConfig() (RepoConfig, error) {
	var conf RepoConfig
	connString, err := GetEnvVar("DB_CONN_STRING")
	if err != nil {
		return conf, err
	}
	return RepoConfig{
		ConnString: connString,
	}, nil
}
