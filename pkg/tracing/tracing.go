package tracing

import (
	"MyMessenger/pkg/logger"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	OTLPEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:""`
}

func Setup(lf fx.Lifecycle, conf *Config, log *logger.Logger, service string) error {
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceName(service))),
	}

	if conf.OTLPEndpoint != "" {
		exporter, err := otlptracehttp.New(context.Background())
		if err != nil {
			return err
		}
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	tp := sdktrace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Warn("otel error", zap.Error(err))
	}))

	lf.Append(fx.StopHook(tp.Shutdown))

	return nil
}
