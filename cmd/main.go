package main

import (
	"fmt"
	"log"

	"github.com/MatthewAraujo/notify/auth"
	"github.com/MatthewAraujo/notify/cmd/api"
	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
)

func main() {

	auth.NewAuth()
	db, err := db.NewMySQLStorage(config.Envs.TursoURl)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
