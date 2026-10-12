package main

import (
	"MyMessenger/services/auth/internal/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(di.GetModule()).Run()
}
