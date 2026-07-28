package di

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/jwt"

	"go.uber.org/fx"
)

var AuthenticatorModule = fx.Options(
	fx.Provide(
		config.GetJwtCheckerConf,
		jwt.NewJwtCheckerConf,
		jwt.NewAuthenticator,
	),
)
