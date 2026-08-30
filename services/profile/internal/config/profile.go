package config

import "MyMessenger/pkg/config"

type ProfileConfig struct {
	ProfileMinNameLen int
	ProfileMaxNameLen int

	ProfileMinBioLen int
	ProfileMaxBioLen int
}

func GetProfileConfig() (*ProfileConfig, error) {
	r := config.NewConfigReader()

	conf := ProfileConfig{
		ProfileMinNameLen: r.GetInt("PROFILE_MIN_NAME_LEN"),
		ProfileMaxNameLen: r.GetInt("PROFILE_MAX_NAME_LEN"),

		ProfileMinBioLen: r.GetInt("PROFILE_MIN_BIO_LEN"),
		ProfileMaxBioLen: r.GetInt("PROFILE_MAX_BIO_LEN"),
	}

	return &conf, r.Err()
}
