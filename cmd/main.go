package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/MatthewAraujo/notify/cmd/api"
)

func main() {
	serve := api.NewAPIServer(fmt.Sprintf(":%d", 8080))
	if err := serve.Start(); err != nil {
		log.Fatalf("error occured while starting server: %s", err)
	}
	slog.Info("server started")
}
