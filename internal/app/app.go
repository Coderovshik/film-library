package app

import (
	"log"

	"github.com/Coderovshik/film-library/internal/config"
	"github.com/Coderovshik/film-library/internal/db"
	"github.com/Coderovshik/film-library/internal/router"
	"github.com/Coderovshik/film-library/internal/user"
)

type App struct {
	Router *router.Router
	Config *config.Config
}

func NewApp(cfg *config.Config) *App {
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewRepository(database.GetDB())
	userService := user.NewService(userRepo, cfg)
	userHandler := user.NewHandler(userService)

	router := router.NewRouter(userHandler)

	return &App{
		Router: router,
		Config: cfg,
	}
}

func (a *App) Run() {
	log.Printf("server running %s", a.Config.Addr())
	if err := a.Router.Run(a.Config.Addr()); err != nil {
		log.Fatal(err)
	}
}
