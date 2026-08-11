package config

import (
	"MyMessenger/pkg/config"
)

type MessageConfig struct {
	MinQuantity int
	MaxQuantity int
}

func GetMessageConfig() (*MessageConfig, error) {
	r := config.NewConfigReader()

	conf := MessageConfig{
		MinQuantity: r.GetInt("MIN_QUANTITY_QUERIES"),
		MaxQuantity: r.GetInt("MAX_QUANTITY_QUERIES"),
	}

	return &conf, r.Err()
}
