package config

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/utils"
	"fmt"
)

func setupProfileValidation(r *config.ConfigReader) {

	getMinMax := func(key1, key2 string) string {
		if r.Err() != nil {
			return ""
		}
		return fmt.Sprintf("required,min=%d,max=%d", r.GetInt(key1), r.GetInt(key2))
	}

	getLenBase64 := func(key string) string {
		if r.Err() != nil {
			return ""
		}
		return fmt.Sprintf("required,len=%d,base64", r.GetInt(key))
	}

	utils.Validator.RegisterAlias("profile_name",
		getMinMax("PROFILE_MIN_NAME_LEN", "PROFILE_MAX_NAME_LEN"),
	)

	utils.Validator.RegisterAlias("profile_username",
		getMinMax("PROFILE_MIN_USERNAME_LEN", "PROFILE_MAX_USERNAME_LEN"),
	)

	utils.Validator.RegisterAlias("profile_pubKey",
		getLenBase64("PROFILE_PUBKEY_LEN"),
	)

	utils.Validator.RegisterAlias("profile_prvKey",
		getLenBase64("PROFILE_PRVKEY_LEN"),
	)

	utils.Validator.RegisterAlias("profile_salt",
		getLenBase64("PROFILE_SALT_LEN"),
	)

	utils.Validator.RegisterAlias("profile_nonce",
		getLenBase64("PROFILE_NONCE_LEN"),
	)
}

func setupMessageValidation(r *config.ConfigReader) {

	getMinMax := func(key1, key2 string) string {
		if r.Err() != nil {
			return ""
		}
		return fmt.Sprintf("required,min=%d,max=%d", r.GetInt(key1), r.GetInt(key2))
	}

	getLenBase64 := func(key string) string {
		if r.Err() != nil {
			return ""
		}
		return fmt.Sprintf("required,len=%d,base64", r.GetInt(key))
	}

	utils.Validator.RegisterAlias("message_text",
		getMinMax("MSG_TEXT_MIN_LEN", "MSG_TEXT_MAX_LEN"),
	)

	utils.Validator.RegisterAlias("message_nonce",
		getLenBase64("MSG_NONCE_LEN"),
	)

	utils.Validator.RegisterAlias("message_key",
		getLenBase64("MSG_KEY_LEN"),
	)
}

func SetupValidator() error {
	r := config.NewConfigReader()

	setupProfileValidation(r)
	setupMessageValidation(r)

	return r.Err()
}
