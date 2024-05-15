package main

import (
	"fmt"
	"log"

	"github.com/MatthewAraujo/notify/cmd/api"
	"github.com/MatthewAraujo/notify/config"
	"github.com/MatthewAraujo/notify/db"
)

func main() {

	url := "libsql://notify-api-matthewaraujo.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicnciLCJpYXQiOjE3MTU3OTE0NDMsImlkIjoiNGZlYjBlZjgtNGJmNS00MjFmLWIxOWMtYjkwYjIyZDNhMjZjIn0.ErOgXIPAAoVxH62MU-x9pff3Fogu_Ej2w9eC4KijNgm0Pj4l8wfdO1TtPsrFWg1ES0CId6P1eOACv8fXfJlADg"
	db, err := db.NewMySQLStorage(url)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
