package logger

import (
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap/zapcore"
)

func NewFxLogger(l *Logger) fxevent.Logger {
	fl := &fxevent.ZapLogger{Logger: l.Logger}
	fl.UseLogLevel(zapcore.DebugLevel)
	fl.UseErrorLevel(zapcore.ErrorLevel)
	return fl
}
