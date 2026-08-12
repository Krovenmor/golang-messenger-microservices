package config

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/utils"
)

func setupProfileValidation(r *config.ConfigValidator) {

	r.SetRangesAlias("profile_name", "PROFILE_MIN_NAME_LEN", "PROFILE_MAX_NAME_LEN")
	r.SetRangesAlias("profile_username", "PROFILE_MIN_USERNAME_LEN", "PROFILE_MAX_USERNAME_LEN")

	r.SetLenBase64Alias("profile_pubKey", "PROFILE_PUBKEY_LEN")
	r.SetLenBase64Alias("profile_prvKey", "PROFILE_PRVKEY_LEN")
	r.SetLenBase64Alias("profile_salt", "PROFILE_SALT_LEN")
	r.SetLenBase64Alias("profile_nonce", "PROFILE_NONCE_LEN")

}

func setupMessageValidation(r *config.ConfigValidator) {

	r.SetRangesBase64Alias("message_text", "MSG_TEXT_MIN_LEN", "MSG_TEXT_MAX_LEN")

	r.SetLenBase64Alias("message_nonce", "MSG_NONCE_LEN")
	r.SetLenBase64Alias("message_key", "MSG_KEY_LEN")

}

func SetupValidator() error {
	r := config.NewConfigValidator(utils.Validator)

	setupProfileValidation(r)
	setupMessageValidation(r)

	return r.Err()
}
