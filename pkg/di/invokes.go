package di

import (
	"MyMessenger/pkg/config"
	"MyMessenger/pkg/logger"
	"MyMessenger/pkg/tracing"
	"context"
	"errors"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const HealthPath = "/healthz"

func InvokeServer(lf fx.Lifecycle, conf config.ServConfig, h http.Handler, log *logger.Logger) {
	instrumented := tracing.HTTPMiddleware(log.HTTPMiddleware(h))
	serv := http.Server{
		Addr: conf.Address,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == HealthPath {
				h.ServeHTTP(w, r)
				return
			}
			instrumented.ServeHTTP(w, r)
		}),
	}
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l, err := net.Listen("tcp", serv.Addr)
			if err != nil {
				return err
			}
			log.Info("server started", zap.String("address", serv.Addr))
			go func() {
				err := serv.Serve(l)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("server stopped with error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping server")
			return serv.Shutdown(ctx)
		},
	})
}
