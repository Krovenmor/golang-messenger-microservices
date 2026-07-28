package main

import (
	"MyMessenger/services/msg/internal/di"
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
