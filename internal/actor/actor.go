package actor

import (
	"context"
	"net/http"
	"time"
)

type Actor struct {
	ID       int32     `json:"id"`
	Name     string    `json:"name"`
	Sex      string    `json:"sex"`
	Birthday time.Time `json:"birthday"`
	Films    []string  `json:"films"`
}

type ActorRepository interface {
	GetActor(ctx context.Context, id int32) (*Actor, error)
	AddActor(ctx context.Context, a *Actor) (*Actor, error)
	DeleteActor(ctx context.Context, id int32) error
	UpdateActor(ctx context.Context, a *Actor) error
	GetActors(ctx context.Context) ([]*Actor, error)
}

type ActorService interface {
	GetActors(ctx context.Context) ([]*ActorResponse, error)
	AddActor(ctx context.Context, req *ActorInfo) (*ActorResponse, error)
	GetActor(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error)
	UpdateActor(ctx context.Context, req *ActorIdInfoRequest) (*ActorResponse, error)
	DeleteActor(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error)
}

type ActorHandler interface {
	GetActors(w http.ResponseWriter, r *http.Request)
	AddActor(w http.ResponseWriter, r *http.Request)
	GetActor(w http.ResponseWriter, r *http.Request)
	UpdateActor(w http.ResponseWriter, r *http.Request)
	DeleteActor(w http.ResponseWriter, r *http.Request)
}

type ActorInfo struct {
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
}

type ActorResponse struct {
	ID    int       `json:"id"`
	Info  ActorInfo `json:"info"`
	Films []string  `json:"films"`
}

type ActorIdRequest struct {
	ID string
}

type ActorIdInfoRequest struct {
	ID   string
	Info ActorInfo
}
