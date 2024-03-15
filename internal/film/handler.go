package film

import (
	"errors"
	"log"
	"net/http"

	"github.com/Coderovshik/film-library/internal/util"
)

var _ FilmHandler = (*Handler)(nil)

type Handler struct {
	service FilmService
}

func NewHandler(fs FilmService) *Handler {
	return &Handler{
		service: fs,
	}
}

func (h *Handler) GetFilms(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetFilms(r.Context(), &GetFilmsRequest{
		SortQuery:  r.URL.Query().Get("sort"),
		FilmQuery:  r.URL.Query().Get("film"),
		ActorQuery: r.URL.Query().Get("actor"),
	})
	if err != nil {
		log.Printf("ERROR: failed to get films err=%s\n", err.Error())

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

func (h *Handler) AddFilm(w http.ResponseWriter, r *http.Request) {
	var req AddFilmRequest
	util.BindJSON(w, r, &req)

	res, err := h.service.AddFilm(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to add film err=%s\n", err.Error())

		var ve *util.ValidationError
		if errors.As(err, &ve) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		if errors.Is(err, ErrFilmActorExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is alreadey bound to the film",
			})
			return
		}

		if errors.Is(err, ErrActorNotExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is non-existent",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) GetFilm(w http.ResponseWriter, r *http.Request) {
	req := FilmIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.GetFilm(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to get film err=%s\n", err.Error())
		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
			util.NotFound(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	var req FilmIdInfoRequest
	util.BindJSON(w, r, &req.Info)
	req.ID = r.PathValue("id")

	res, err := h.service.UpdateFilm(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to update film err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
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

		if errors.Is(err, ErrEmptyUpdate) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeValidation,
				Body:      "empty update",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	req := FilmIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.DeleteFilm(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to delete film err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
			util.NotFound(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) GetFilmActors(w http.ResponseWriter, r *http.Request) {
	req := FilmIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.GetFilmActors(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to get film related actors err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
			util.NotFound(w, r)
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) AddFilmActors(w http.ResponseWriter, r *http.Request) {
	var req FilmActorsRequest
	util.BindJSON(w, r, &req.ActorIDs)
	req.ID = r.PathValue("id")

	res, err := h.service.AddFilmActors(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to add film related actors err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
			util.NotFound(w, r)
			return
		}

		if errors.Is(err, ErrFilmActorExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is alreadey bound to the film",
			})
			return
		}

		if errors.Is(err, ErrActorNotExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is non-existent",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) DeleteFilmActors(w http.ResponseWriter, r *http.Request) {
	var req FilmActorsRequest
	util.BindJSON(w, r, &req.ActorIDs)
	req.ID = r.PathValue("id")

	res, err := h.service.DeleteFilmActors(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to add film related actors err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrFilmNotExist) {
			util.NotFound(w, r)
			return
		}

		if errors.Is(err, ErrFilmActorExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is alreadey bound to the film",
			})
			return
		}

		if errors.Is(err, ErrActorNotExist) {
			util.JSON(w, r, http.StatusBadRequest, &util.ErrorMessage{
				ErrorType: util.ErrorTypeConflict,
				Body:      "one of the provided actors is non-existent",
			})
			return
		}

		util.InternalServerError(w, r)
		return
	}

	util.JSON(w, r, http.StatusOK, res)
}
