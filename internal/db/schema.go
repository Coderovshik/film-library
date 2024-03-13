package db

import (
	"embed"
	"log"

	"github.com/Coderovshik/film-library/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var schemaFS embed.FS

func NewMigrator(cfg *config.Config) *migrate.Migrate {
	d, err := iofs.New(schemaFS, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	return m
}
