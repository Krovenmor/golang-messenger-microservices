package config

import "MyMessenger/pkg/config"

type MessageConfig struct {
	MinNameLen int
	MaxNameLen int

	MinUserNameLen int
	MaxUserNameLen int

	MinMsgLen int
	MaxMsgLen int

	MinKeysLen   int
	MaxPubKeyLen int
	MaxPrvKeyLen int
	MaxSaltLen   int

	MinQuantity int
	MaxQuantity int
}

func GetMessageConfig() (*MessageConfig, error) {
	var (
		conf MessageConfig
		err  error
	)

	get := func(envKey string) int {
		if err != nil {
			return 0
		}
		var val int
		val, err = config.GetEnvVarInt(envKey)
		return val
	}

	conf = MessageConfig{
		MinNameLen:     get("MIN_NAME_LEN"),
		MaxNameLen:     get("MAX_NAME_LEN"),
		MinUserNameLen: get("MIN_USERNAME_LEN"),
		MaxUserNameLen: get("MAX_USERNAME_LEN"),
		MinMsgLen:      get("MIN_MSG_LEN"),
		MaxMsgLen:      get("MAX_MSG_LEN"),
		MinKeysLen:     get("MIN_KEYS_LEN"),
		MaxPubKeyLen:   get("MAX_PUBKEY_LEN"),
		MaxPrvKeyLen:   get("MAX_PRVKEY_LEN"),
		MaxSaltLen:     get("MAX_SALT_LEN"),
		MinQuantity:    get("MIN_QUANTITY_QUERIES"),
		MaxQuantity:    get("MAX_QUANTITY_QUERIES"),
	}

	if err != nil {
		return nil, err
	}

	return &conf, nil
}
