package db

import (
	"database/sql"
	"log"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewMySQLStorage(url string) (*sql.DB, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return db, nil
}
