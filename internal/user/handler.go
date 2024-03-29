package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/Coderovshik/film-library/internal/util"
)

type Handler struct {
	service UserService
}

func NewHandler(s UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if ok := util.BindJSON(w, r, &req); !ok {
		return
	}

	res, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to create user err=%s\n", err.Error())

		var ve *util.ValidationError
		if errors.As(err, &ve) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		if errors.Is(err, ErrUserExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "User already exists",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if ok := util.BindJSON(w, r, &req); !ok {
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to login user err=%s\n", err.Error())

		if errors.Is(err, ErrUserNotExist) || errors.Is(err, ErrPasswordIncorrect) {
			util.JSON(w, r, http.StatusUnauthorized, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      "incorrect password",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.SetJWTCookie(w, res.AccessToken)
	util.OK(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt")
	if err != nil {
		log.Printf("ERROR: logout failed err=%s", err.Error())
		if errors.Is(err, http.ErrNoCookie) {
			util.Unauthorized(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.UnsetJWTCookie(w)
	util.OK(w, r)
}
