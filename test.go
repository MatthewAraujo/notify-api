package main

import (
	"log"

	"github.com/MatthewAraujo/notify/encrypt"
)

func main() {
	token := "ghs_JpoTODOb6M6iZveqBZvhzbSUMMxqRu0bV040"

	token, err := encrypt.EncryptToken(token)
	if err != nil {
		panic(err)
	}

	log.Printf("Decrypted token: %s", token)

}
