package di

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/logger"
	"MyMessenger/pkg/tracing"

	"go.uber.org/fx"
)

func LoggerModule(service string) fx.Option {
	return fx.Options(
		fx.Provide(
			config.Parse[logger.Config],
			func(lf fx.Lifecycle, conf *logger.Config) (*logger.Logger, error) {
				return logger.New(lf, conf, service)
			},
		),
		fx.WithLogger(logger.NewFxLogger),
	)
}

func TracingModule(service string) fx.Option {
	return fx.Options(
		fx.Provide(
			config.Parse[tracing.Config],
		),
		fx.Invoke(
			func(lf fx.Lifecycle, conf *tracing.Config, log *logger.Logger) error {
				return tracing.Setup(lf, conf, log, service)
			},
		),
	)
}
