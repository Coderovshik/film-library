package main

import (
	"log"

	"github.com/Coderovshik/film-library/internal/config"
	"github.com/Coderovshik/film-library/internal/db"
)

func main() {
	cfg := config.New()
	m := db.NewMigrator(cfg)

	log.Println("migrating")
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
