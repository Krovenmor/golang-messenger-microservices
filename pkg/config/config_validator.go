package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ConfigValidator struct {
	*ConfigReader
	*validator.Validate
}

func NewConfigValidator(v *validator.Validate) *ConfigValidator {
	return &ConfigValidator{NewConfigReader(), v}
}

func (v *ConfigValidator) SetRangesAlias(alias, keyFrom, keyTo string) {
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,min=%d,max=%d", v.GetInt(keyFrom), v.GetInt(keyTo)),
	)
}

func (v *ConfigValidator) SetLenBase64Alias(alias, keyLen string) {
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,len=%d,base64", v.GetInt(keyLen)),
	)
}
