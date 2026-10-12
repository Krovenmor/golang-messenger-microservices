package config

import (
	"MyMessenger/pkg/utils"
	"fmt"

	"github.com/go-playground/validator/v10/non-standard/validators"
)

func SetupValidator(conf *ValidationConfig) error {
	err := utils.Validator.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return err
	}
	utils.Validator.RegisterAlias("auth_password", fmt.Sprintf("required,min=%d,max=%d,notblank", conf.PassMinLength, conf.PassMaxLength))
	utils.Validator.RegisterAlias("auth_login", fmt.Sprintf("required,min=%d,max=%d,alphanum", conf.LoginMinLength, conf.LoginMaxLength))
	utils.Validator.RegisterAlias("refresh_token", fmt.Sprintf("required,len=%d", conf.RefreshTokenLen))
	return nil
}
