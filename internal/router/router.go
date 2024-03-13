package router

import (
	"net/http"

	"github.com/Coderovshik/film-library/internal/user"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(uh user.UserHandler) *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup/", uh.CreateUser)
	mux.HandleFunc("POST /signin/", uh.Login)

	return &Router{
		mux: mux,
	}
}

func (r *Router) Run(addr string) {
	http.ListenAndServe(addr, r.mux)
}
