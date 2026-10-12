package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`
}

type Logger struct {
	*zap.Logger
}

func New(lf fx.Lifecycle, conf *Config, service string) (*Logger, error) {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(conf.Level)
	zc.OutputPaths = []string{"stdout"}
	zc.EncoderConfig.TimeKey = "time"
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zc.InitialFields = map[string]any{"service": service}

	l, err := zc.Build()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(l)

	lf.Append(fx.StopHook(func() {
		_ = l.Sync()
	}))

	return &Logger{l}, nil
}

func (l *Logger) Ctx(ctx context.Context) *zap.Logger {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return l.Logger
	}
	return l.With(
		zap.Stringer("trace_id", sc.TraceID()),
		zap.Stringer("span_id", sc.SpanID()),
	)
}
