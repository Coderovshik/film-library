package user

import (
	"encoding/json"
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
	log.Printf("INFO: incoming reqeust path=%s from=%s\n", r.URL.Path, r.RemoteAddr)

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body err=%s\n", err.Error())
		util.InternalServerError(w, r)
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

		// TODO: add check for ErrUserExist

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: incoming reqeust path=%s from=%s\n", r.URL.Path, r.RemoteAddr)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body err=%s\n", err.Error())
		util.InternalServerError(w, r)
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to create user err=%s\n", err.Error())

		// TODO: add check for ErrUserNotExist and Password comparison failure

		util.InternalServerError(w, r)
		return
	}

	util.SetJWTCookie(w, res.AccessToken)
	util.OK(w, r)
}
