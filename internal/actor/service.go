package actor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Coderovshik/film-library/internal/config"
)

var (
	ErrIdInvalid = errors.New("invalid id")
)

var _ ActorService = (*Service)(nil)

type Service struct {
	repo ActorRepository
}

func NewService(ar ActorRepository, cfg *config.Config) *Service {
	return &Service{
		repo: ar,
	}
}

func (s *Service) GetActors(ctx context.Context) ([]*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	actors, err := s.repo.GetActors(ctx)
	if err != nil {
		log.Printf("ERROR: failed to get actor records from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := make([]*ActorResponse, 0, len(actors))
	for _, v := range actors {
		res = append(res, ToActorResponse(v))
	}

	return res, nil
}

func (s *Service) AddActor(ctx context.Context, req *ActorInfo) (*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	vErr := ValidateEmptyActorInfo(req)
	if vErr != nil {
		log.Printf("ERROR: failed request empty validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	vErr = ValidateFormatActorInfo(req)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}
	actor := ToActor(req)

	actor, err := s.repo.AddActor(ctx, actor)
	if err != nil {
		log.Printf("ERROR: failed to create actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) GetActor(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error) {
	const op = "actor.Service.GetActor"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actor, err := s.repo.GetActor(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) UpdateActor(ctx context.Context, req *ActorIdInfoRequest) (*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	vErr := ValidateFormatActorInfo(&req.Info)
	if vErr != nil {
		log.Printf("ERROR: failed request format validation\n")
		return nil, fmt.Errorf("%s: %w", op, vErr)
	}

	actor := ToActor(&req.Info)
	actor.ID = int32(id)

	err = s.repo.UpdateActor(ctx, actor)
	if err != nil {
		log.Printf("ERROR: failed to update actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	actor, err = s.repo.GetActor(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) DeleteActor(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	actor, err := s.repo.GetActor(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to get actor record from repository\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.repo.DeleteActor(ctx, int32(id))
	if err != nil {
		log.Printf("ERROR: failed to delete actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}
