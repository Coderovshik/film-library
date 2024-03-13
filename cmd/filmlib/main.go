package main

import (
	"github.com/Coderovshik/film-library/internal/app"
	"github.com/Coderovshik/film-library/internal/config"
)

func main() {
	cfg := config.New()
	app := app.NewApp(cfg)

	app.Run()
}
