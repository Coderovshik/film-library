package router

import (
	"net/http"

	"github.com/Coderovshik/film-library/internal/actor"
	"github.com/Coderovshik/film-library/internal/config"
	"github.com/Coderovshik/film-library/internal/film"
	"github.com/Coderovshik/film-library/internal/middleware"
	"github.com/Coderovshik/film-library/internal/user"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(cfg *config.Config, uh user.UserHandler, ah actor.ActorHandler, fh film.FilmHandler) *Router {
	mux := http.NewServeMux()

	authMW := middleware.NewAuthMiddleware(cfg.SigningKey, false)
	adminOnlyMW := middleware.NewAuthMiddleware(cfg.SigningKey, true)
	logMW := middleware.NewLogMiddleware()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.Handle("POST /signup", logMW(http.HandlerFunc(uh.CreateUser)))
	mux.Handle("POST /signin", logMW(http.HandlerFunc(uh.Login)))
	mux.Handle("DELETE /signout", logMW(http.HandlerFunc(uh.Logout)))

	mux.Handle("GET /actors", logMW(authMW(http.HandlerFunc(ah.GetActors))))
	mux.Handle("POST /actors", logMW(adminOnlyMW(http.HandlerFunc(ah.AddActor))))
	mux.Handle("GET /actors/{id}", logMW(authMW(http.HandlerFunc(ah.GetActor))))
	mux.Handle("PUT /actors/{id}", logMW(adminOnlyMW(http.HandlerFunc(ah.UpdateActor))))
	mux.Handle("DELETE /actors/{id}", logMW(adminOnlyMW(http.HandlerFunc(ah.DeleteActor))))

	mux.Handle("GET /films", logMW(authMW(http.HandlerFunc(fh.GetFilms))))
	mux.Handle("POST /films", logMW(adminOnlyMW(http.HandlerFunc(fh.AddFilm))))
	mux.Handle("GET /films/{id}", logMW(authMW(http.HandlerFunc(fh.GetFilm))))
	mux.Handle("PUT /films/{id}", logMW(adminOnlyMW(http.HandlerFunc(fh.UpdateFilm))))
	mux.Handle("DELETE /films/{id}", logMW(adminOnlyMW(http.HandlerFunc(fh.DeleteFilm))))
	mux.Handle("GET /films/{id}/actors", logMW(authMW(http.HandlerFunc(fh.GetFilmActors))))
	mux.Handle("PUT /films/{id}/actors", logMW(adminOnlyMW(http.HandlerFunc(fh.AddFilmActors))))
	mux.Handle("DELETE /films/{id}/actors", logMW(adminOnlyMW(http.HandlerFunc(fh.DeleteFilmActors))))

	return &Router{
		mux: mux,
	}
}

func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}
