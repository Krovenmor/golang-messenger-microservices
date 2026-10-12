package config

import (
	"reflect"

	"github.com/caarlos0/env/v11"
)

var envParsers = map[reflect.Type]env.ParserFunc{
	reflect.TypeFor[[]byte](): func(v string) (any, error) {
		return []byte(v), nil
	},
}

func Parse[T any]() (*T, error) {
	conf, err := env.ParseAsWithOptions[T](env.Options{
		RequiredIfNoDef: true,
		FuncMap:         envParsers,
	})
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
