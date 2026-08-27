package config

import "MyMessenger/pkg/config"

type MediaSaverConfig struct {
	PublicPhotosSavingPath string
}

type HandlerConfig struct {
	AvatarsSizeLimit int
}

func GetMediaSaverConfig() (*MediaSaverConfig, error) {
	r := config.NewConfigReader()

	conf := MediaSaverConfig{
		PublicPhotosSavingPath: r.GetString("PUBLIC_PHOTOS_PATH"),
	}

	return &conf, r.Err()
}

func GetHandlerConfig() (*HandlerConfig, error) {
	r := config.NewConfigReader()

	conf := HandlerConfig{
		AvatarsSizeLimit: r.GetInt("AVATARS_SIZE_LIMIT"),
	}

	return &conf, r.Err()
}
