package main

import (
	"fmt"
	"log/slog"

	"github.com/MatthewAraujo/notify/cmd/api"
)

func main() {
	server := api.NewAPIServer(fmt.Sprintf(":%d", 8080))
	if err := server.Start(); err != nil {
		fmt.Printf("error occured while starting server: %s", err)
	}
	slog.Info("server started")

}
