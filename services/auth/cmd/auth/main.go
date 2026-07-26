package main

import (
	"MyMessenger/services/auth/internal/di"
	"log"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(di.GetModule())
	if app == nil {
		log.Fatalf("app == nil")
	}
	app.Run()
}
