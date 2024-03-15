package app

import (
	"log"

	"github.com/Coderovshik/film-library/internal/actor"
	"github.com/Coderovshik/film-library/internal/config"
	"github.com/Coderovshik/film-library/internal/db"
	"github.com/Coderovshik/film-library/internal/film"
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

	actorRepo := actor.NewRepository(database.GetDB())
	actorService := actor.NewService(actorRepo, cfg)
	actorHandler := actor.NewHandler(actorService)

	filmRepo := film.NewRepository(database.GetDB())
	filmService := film.NewService(filmRepo)
	filmHandler := film.NewHandler(filmService)

	router := router.NewRouter(cfg, userHandler, actorHandler, filmHandler)

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
