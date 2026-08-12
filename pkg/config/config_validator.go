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
	if v.Err() != nil {
		return
	}
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,alphanum,min=%d,max=%d", v.GetInt(keyFrom), v.GetInt(keyTo)),
	)
}

func (v *ConfigValidator) SetRangesBase64Alias(alias, keyFrom, keyTo string) {
	if v.Err() != nil {
		return
	}
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,min=%d,max=%d,base64", v.GetInt(keyFrom), v.GetInt(keyTo)),
	)
}

func (v *ConfigValidator) SetLenBase64Alias(alias, keyLen string) {
	if v.Err() != nil {
		return
	}
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,len=%d,base64", v.GetInt(keyLen)),
	)
}

func (v *ConfigValidator) SetLenAlias(alias, keyLen string) {
	if v.Err() != nil {
		return
	}
	v.RegisterAlias(
		alias,
		fmt.Sprintf("required,alphanum,len=%d", v.GetInt(keyLen)),
	)
}
