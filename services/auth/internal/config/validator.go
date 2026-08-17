package config

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/utils"
)

func SetupValidator() error {
	r := config.NewConfigValidator(utils.Validator)

	r.SetRangesAlias("auth_password", "PASS_MIN_LENGTH", "PASS_MAX_LENGTH")
	r.SetRangesAlias("auth_login", "LOGIN_MIN_LENGTH", "LOGIN_MAX_LENGTH")
	r.SetLenAliasBasic("refresh_token", "REFRESH_TOKEN_LEN", "")

	return r.Err()
}
