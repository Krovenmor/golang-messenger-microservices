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

func (v *ConfigValidator) SetRangesAliasBasic(alias, keyFrom, keyTo, ext string) {
	if v.Err() != nil {
		return
	}
	aliasStr := fmt.Sprintf("required,min=%d,max=%d", v.GetInt(keyFrom), v.GetInt(keyTo))
	if len(ext) > 0 {
		aliasStr += "," + ext
	}
	v.RegisterAlias(alias, aliasStr)
}

func (v *ConfigValidator) SetLenAliasBasic(alias, keyLen, ext string) {
	if v.Err() != nil {
		return
	}
	aliasStr := fmt.Sprintf("required,len=%d", v.GetInt(keyLen))
	if len(ext) > 0 {
		aliasStr += "," + ext
	}
	v.RegisterAlias(alias, aliasStr)
}

func (v *ConfigValidator) SetRangesAlias(alias, keyFrom, keyTo string) {
	v.SetRangesAliasBasic(alias, keyFrom, keyTo, "alphanum")
}

func (v *ConfigValidator) SetRangesBase64Alias(alias, keyFrom, keyTo string) {
	v.SetRangesAliasBasic(alias, keyFrom, keyTo, "base64")
}

func (v *ConfigValidator) SetLenBase64Alias(alias, keyLen string) {
	v.SetLenAliasBasic(alias, keyLen, "base64")
}

func (v *ConfigValidator) SetLenAlias(alias, keyLen string) {
	v.SetLenAliasBasic(alias, keyLen, "alphanum")
}
