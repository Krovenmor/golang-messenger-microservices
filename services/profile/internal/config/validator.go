package config

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/utils"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func registerCustomUTFValidator(v *validator.Validate) error {
	return v.RegisterValidation("valid_utf8_printable", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		for _, r := range str {
			if unicode.IsSpace(r) {
				return false
			}
			if unicode.IsControl(r) {
				return false
			}
			if r == unicode.ReplacementChar {
				return false
			}
		}
		return true
	})
}

func setupProfileValidation(r *config.ConfigValidator) {

	r.SetRangesAliasBasic("profile_name", "PROFILE_MIN_NAME_LEN", "PROFILE_MAX_NAME_LEN", "valid_utf8_printable")
	r.SetRangesAlias("profile_username", "PROFILE_MIN_USERNAME_LEN", "PROFILE_MAX_USERNAME_LEN")

	r.SetLenBase64Alias("profile_pubKey", "PROFILE_PUBKEY_LEN")
	r.SetLenBase64Alias("profile_prvKey", "PROFILE_PRVKEY_LEN")
	r.SetLenBase64Alias("profile_salt", "PROFILE_SALT_LEN")
	r.SetLenBase64Alias("profile_nonce", "PROFILE_NONCE_LEN")

}

func SetupValidator() error {
	r := config.NewConfigValidator(utils.Validator)

	registerCustomUTFValidator(utils.Validator)
	setupProfileValidation(r)

	return r.Err()
}
