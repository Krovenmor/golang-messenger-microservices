package di

import (
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra"
	"MyMessenger/services/auth/internal/service"
	web "MyMessenger/services/auth/internal/transport/http"
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

func GetModule() fx.Option {
	return fx.Options(

		// Configs
		fx.Provide(
			config.GetAuthConfig,
			config.GetRepoConfig,
			config.GetServConfig,
		),

		// Pool
		fx.Provide(
			infra.NewPool,
		),

		// Make migrations
		fx.Invoke(
			infra.MakeMigrations,
		),

		// Repo Service
		fx.Provide(
			fx.Annotate(
				infra.NewRepo,
				fx.As(new(service.AuthRepo)),
			),
		),

		// Auth Service
		fx.Provide(
			fx.Annotate(
				service.NewAuth,
				fx.As(new(service.AuthService)),
			),
		),

		// ServeMux
		fx.Provide(
			http.NewServeMux,
		),

		fx.Provide(
			web.NewHandler,
		),

		fx.Provide(
			(*web.Handler).RegisterRoutes,
		),

		// Final, main invoke
		fx.Invoke(
			InvokeServer,
		),
	)
}
