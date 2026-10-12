package http

import (
	"MyMessenger/pkg/httpx"
	"MyMessenger/services/auth/internal/config"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func NewValidator(conf *config.ValidationConfig) (*validator.Validate, error) {
	v := httpx.NewValidator()
	err := v.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return nil, err
	}
	v.RegisterAlias("auth_password", fmt.Sprintf("required,min=%d,max=%d,notblank", conf.PassMinLength, conf.PassMaxLength))
	v.RegisterAlias("auth_login", fmt.Sprintf("required,min=%d,max=%d,alphanum", conf.LoginMinLength, conf.LoginMaxLength))
	v.RegisterAlias("refresh_token", fmt.Sprintf("required,len=%d", conf.RefreshTokenLen))
	return v, nil
}
