package film

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Coderovshik/film-library/internal/util"
)

var (
	ErrIdInvalid = errors.New("invalid id")
)

type Service struct {
	repo FilmRepository
}

func NewService(fr FilmRepository) *Service {
	return &Service{
		repo: fr,
	}
}

func (s *Service) GetFilms(ctx context.Context, req *GetFilmsRequest) ([]*FilmResponse, error) {
	const op = "film.Service.GetFilms"

	vErr := ValidateGetFilmsRequest(req)
	if vErr != nil {
		log.Printf("ERROR: failed request validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	q := ToQuery(req)

	films, err := s.repo.GetFilms(ctx, q)
	if err != nil {
		log.Printf("ERROR: failed to get films")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]*FilmResponse, 0, len(films))
	for _, v := range films {
		res = append(res, ToFilmResponse(v))
	}

	return res, nil
}

func (s *Service) AddFilm(ctx context.Context, req *AddFilmRequest) (*FilmResponse, error) {
	const op = "film.Service.AddFilm"

	vErr := ValidateEmptyFilmInfo(&req.Info)
	if vErr != nil {
		log.Printf("ERROR: failed request empty validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	vErr = ValidateFormatFilmInfo(&req.Info, false)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	film := ToFilm(&req.Info)

	film, err := s.repo.AddFilm(ctx, film)
	if err != nil {
		log.Printf("ERROR: failed to add film information\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fa := &FilmActors{
		ID:       film.ID,
		ActorIDs: ToActorIDs32(util.RemoveDuplicateInt(req.ActorIDs)),
	}
	err = s.repo.AddFilmActors(ctx, fa)
	if err != nil {
		log.Printf("ERROR: failed to bind provided actors and film\n")
		s.repo.DeleteFilm(ctx, film.ID)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	film, err = s.repo.GetFilm(ctx, film.ID)
	if err != nil {
		log.Printf("ERROR: failed to get added film information\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToFilmResponse(film)

	return res, nil
}

func (s *Service) GetFilm(ctx context.Context, req *FilmIdRequest) (*FilmResponse, error) {
	const op = "film.Service.GetFilm"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actor, err := s.repo.GetFilm(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToFilmResponse(actor)

	return res, nil
}

func (s *Service) UpdateFilm(ctx context.Context, req *FilmIdInfoRequest) (*FilmResponse, error) {
	const op = "film.Service.UpdateFilm"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	vErr := ValidateFormatFilmInfo(&req.Info, true)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}

	film := ToFilm(&req.Info)
	film.ID = int32(id)

	err = s.repo.UpdateFilm(ctx, film)
	if err != nil {
		log.Printf("ERROR: failed to update film record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	film, err = s.repo.GetFilm(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get film record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToFilmResponse(film)

	return res, nil
}

func (s *Service) DeleteFilm(ctx context.Context, req *FilmIdRequest) (*FilmResponse, error) {
	const op = "film.Service.DeleteFilm"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	film, err := s.repo.GetFilm(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get film record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.repo.DeleteFilm(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to delete film record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToFilmResponse(film)

	return res, nil
}

func (s *Service) GetFilmActors(ctx context.Context, req *FilmIdRequest) ([]*ActorShortResponse, error) {
	const op = "film.Service.GetFilmActors"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actors, err := s.repo.GetFilmActors(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get film actors from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(actors) == 0 {
		log.Printf("ERROR: no actors found")
		return nil, fmt.Errorf("%s: %w", op, ErrZeroActors)
	}

	res := ToActorsShortRespose(actors)

	return res, nil
}

func (s *Service) AddFilmActors(ctx context.Context, req *FilmActorsRequest) ([]*ActorShortResponse, error) {
	const op = "film.Service.AddFilmActors"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	if req.ActorIDs == nil || len(req.ActorIDs) == 0 {
		log.Printf("ERROR: empty update\n")
		return nil, fmt.Errorf("%s: %w", op, ErrEmptyUpdate)
	}

	fa := &FilmActors{
		ID:       int32(id),
		ActorIDs: ToActorIDs32(util.RemoveDuplicateInt(req.ActorIDs)),
	}
	err = s.repo.AddFilmActors(ctx, fa)
	if err != nil {
		log.Printf("ERROR: failed to bind provided actors and film\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	actors, err := s.repo.GetFilmActors(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get film actors from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorsShortRespose(actors)

	return res, nil
}

func (s *Service) DeleteFilmActors(ctx context.Context, req *FilmActorsRequest) ([]*ActorShortResponse, error) {
	const op = "film.Service.AddFilmActors"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	if req.ActorIDs == nil || len(req.ActorIDs) == 0 {
		log.Printf("ERROR: empty update\n")
		return nil, fmt.Errorf("%s: %w", op, ErrEmptyUpdate)
	}

	fa := &FilmActors{
		ID:       int32(id),
		ActorIDs: ToActorIDs32(util.RemoveDuplicateInt(req.ActorIDs)),
	}
	err = s.repo.DeleteFilmActors(ctx, fa)
	if err != nil {
		log.Printf("ERROR: failed to bind provided actors and film\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	actors, err := s.repo.GetFilmActors(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get film actors from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorsShortRespose(actors)

	return res, nil
}
