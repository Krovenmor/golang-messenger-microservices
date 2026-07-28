package di

import (
	"MyMessenger/pkg/config"
	"context"
	"log"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func InvokeServer(lf fx.Lifecycle, conf config.ServConfig, h http.Handler) {
	serv := http.Server{
		Addr:    conf.Address,
		Handler: h,
	}
	lf.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l, err := net.Listen("tcp", serv.Addr)
			if err != nil {
				return err
			}
			log.Printf("Server started on %q", serv.Addr)
			go serv.Serve(l)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Stopping server")
			return serv.Shutdown(ctx)
		},
	})
}
