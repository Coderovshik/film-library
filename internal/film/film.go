package film

import (
	"context"
	"net/http"
	"time"
)

type Film struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"releasedate"`
	Rating      int32     `json:"rating"`
	Actors      []string  `json:"actors"`
}

type FilmRepository interface {
	GetFilm(ctx context.Context, id int32) (*Film, error)
	AddFilm(ctx context.Context, f *Film) (*Film, error)
	DeleteFilm(ctx context.Context, id int32) error
	UpdateFilm(ctx context.Context, f *Film) error
	GetFilms(ctx context.Context, q *Query) ([]*Film, error)
	GetFilmActors(ctx context.Context, id int32) ([]*ActorShort, error)
	AddFilmActors(ctx context.Context, fa *FilmActors) error
	DeleteFilmActors(ctx context.Context, fa *FilmActors) error
}

type FilmService interface {
	GetFilms(ctx context.Context, req *GetFilmsRequest) ([]*FilmResponse, error)
	AddFilm(ctx context.Context, req *AddFilmRequest) (*FilmResponse, error)
	GetFilm(ctx context.Context, req *FilmIdRequest) (*FilmResponse, error)
	UpdateFilm(ctx context.Context, req *FilmIdInfoRequest) (*FilmResponse, error)
	DeleteFilm(ctx context.Context, req *FilmIdRequest) (*FilmResponse, error)
	GetFilmActors(ctx context.Context, req *FilmIdRequest) ([]*ActorShortResponse, error)
	AddFilmActors(ctx context.Context, req *FilmActorsRequest) ([]*ActorShortResponse, error)
	DeleteFilmActors(ctx context.Context, req *FilmActorsRequest) ([]*ActorShortResponse, error)
}

type FilmHandler interface {
	GetFilms(w http.ResponseWriter, r *http.Request)
	AddFilm(w http.ResponseWriter, r *http.Request)
	GetFilm(w http.ResponseWriter, r *http.Request)
	UpdateFilm(w http.ResponseWriter, r *http.Request)
	DeleteFilm(w http.ResponseWriter, r *http.Request)
	GetFilmActors(w http.ResponseWriter, r *http.Request)
	AddFilmActors(w http.ResponseWriter, r *http.Request)
	DeleteFilmActors(w http.ResponseWriter, r *http.Request)
}

type Filmhandler interface{}

type Query struct {
	Sort  []string
	Actor string
	Film  string
}

type FilmActors struct {
	ID       int32
	ActorIDs []int32
}

type ActorShort struct {
	ID   int32
	Name string
}

type GetFilmsRequest struct {
	SortQuery  string
	FilmQuery  string
	ActorQuery string
}

type AddFilmRequest struct {
	Info     FilmInfo `json:"info"`
	ActorIDs []int    `json:"actorIds"`
}

type FilmInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ReleaseDate string `json:"releasedate"`
	Rating      int    `json:"rating"`
}

type FilmResponse struct {
	ID     int      `json:"id"`
	Info   FilmInfo `json:"info"`
	Actors []string `json:"actors,omitempty"`
}

type FilmIdRequest struct {
	ID string
}

type FilmIdInfoRequest struct {
	ID   string
	Info FilmInfo
}

type FilmActorsRequest struct {
	ID       string
	ActorIDs []int
}

type ActorShortResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
