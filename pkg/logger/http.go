package logger

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *Logger) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)

		level := zapcore.InfoLevel
		if m.Code >= http.StatusInternalServerError {
			level = zapcore.ErrorLevel
		}

		l.Ctx(r.Context()).Log(level, "http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("route", r.Pattern),
			zap.Int("status", m.Code),
			zap.Duration("duration", m.Duration),
			zap.Int64("bytes", m.Written),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}
