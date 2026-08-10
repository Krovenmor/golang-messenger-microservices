package config

import "MyMessenger/pkg/config"

type MessageConfig struct {
	MinNameLen int
	MaxNameLen int

	MinUserNameLen int
	MaxUserNameLen int

	MinMsgLen   int
	MaxMsgLen   int
	MsgNonceLen int
	MsgKeysLen  int

	MinKeysLen   int
	MaxPubKeyLen int
	MaxPrvKeyLen int
	MaxSaltLen   int

	MinNonceLen int
	MaxNonceLen int

	MinQuantity int
	MaxQuantity int
}

func GetMessageConfig() (*MessageConfig, error) {
	r := config.NewConfigReader()

	conf := MessageConfig{
		MinNameLen:     r.GetInt("MIN_NAME_LEN"),
		MaxNameLen:     r.GetInt("MAX_NAME_LEN"),
		MinUserNameLen: r.GetInt("MIN_USERNAME_LEN"),
		MaxUserNameLen: r.GetInt("MAX_USERNAME_LEN"),

		MinMsgLen:   r.GetInt("MIN_MSG_LEN"),
		MaxMsgLen:   r.GetInt("MAX_MSG_LEN"),
		MsgNonceLen: r.GetInt("MSG_NONCE_LEN"),
		MsgKeysLen:  r.GetInt("MSG_KEYS_LEN"),

		MinKeysLen:   r.GetInt("MIN_KEYS_LEN"),
		MaxPubKeyLen: r.GetInt("MAX_PUBKEY_LEN"),
		MaxPrvKeyLen: r.GetInt("MAX_PRVKEY_LEN"),
		MaxSaltLen:   r.GetInt("MAX_SALT_LEN"),
		MinQuantity:  r.GetInt("MIN_QUANTITY_QUERIES"),
		MaxQuantity:  r.GetInt("MAX_QUANTITY_QUERIES"),

		MinNonceLen: r.GetInt("MIN_NONCE_LEN"),
		MaxNonceLen: r.GetInt("MAX_NONCE_LEN"),
	}

	return &conf, r.Err()
}
