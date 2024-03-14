package actor

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Coderovshik/film-library/internal/util"
)

type Handler struct {
	service ActorService
}

func NewHandler(as ActorService) *Handler {
	return &Handler{
		service: as,
	}
}

func (h *Handler) GetActors(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetActors(r.Context())
	if err != nil {
		log.Printf("ERROR: failed to get actors err=%s\n", err.Error())
		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) AddActor(w http.ResponseWriter, r *http.Request) {
	var req ActorInfo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: failed to decode request body err=%s\n", err.Error())
		util.InternalServerError(w, r)
		return
	}

	res, err := h.service.AddActor(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to add actor err=%s\n", err.Error())

		var ve *util.ValidationError
		if errors.As(err, &ve) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) GetActor(w http.ResponseWriter, r *http.Request) {
	req := ActorIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.GetActor(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to get actor err=%s\n", err.Error())
		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			util.NotFound(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	req := ActorIdInfoRequest{
		ID: r.PathValue("id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&req.Info); err != nil {
		log.Printf("ERROR: failed to decode request body err=%s\n", err.Error())
		util.InternalServerError(w, r)
		return
	}

	res, err := h.service.UpdateActor(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to update actor err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			util.NotFound(w, r)
			return
		}

		var ve *util.ValidationError
		if errors.As(err, &ve) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	req := ActorIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.DeleteActor(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to get actor err=%s\n", err.Error())
		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			util.NotFound(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}
