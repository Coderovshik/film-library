package user

import (
	"context"
	"net/http"
)

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Passhash string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	accessToken string
}
