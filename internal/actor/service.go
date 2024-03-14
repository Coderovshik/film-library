package actor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Coderovshik/film-library/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrPermissionDenied = errors.New("no necessary permissions")
	ErrTokenInvalid     = errors.New("invalid api token")
	ErrIdInvalid        = errors.New("invalid id")
)

var _ ActorService = (*Service)(nil)

type Service struct {
	repo       ActorRepository
	signingKey string
}

func (s *Service) GetActors(ctx context.Context) ([]*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	token := ctx.Value(util.JWTContextKey("jwt")).(string)
	uc, err := util.ParseUserClaims(token, s.signingKey)
	if err != nil {
		log.Printf("ERROR: bad token\n")

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Printf("ERROR: jwt signature is invalid\n")
		}
		if errors.Is(err, util.ErrUnknownClaimsType) {
			log.Printf("ERROR: unknown token claims\n")
		}

		return nil, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	if !uc.IsAdmin {
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

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

	token := ctx.Value(util.JWTContextKey("jwt")).(string)
	uc, err := util.ParseUserClaims(token, s.signingKey)
	if err != nil {
		log.Printf("ERROR: bad token\n")

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Printf("ERROR: jwt signature is invalid\n")
		}
		if errors.Is(err, util.ErrUnknownClaimsType) {
			log.Printf("ERROR: unknown token claims\n")
		}

		return nil, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	if !uc.IsAdmin {
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	err = ValidateActorInfo(req)
	if err != nil {
		log.Printf("ERROR: failed request validation\n")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	actor := ToActor(req)

	actor, err = s.repo.AddActor(ctx, actor)
	if err != nil {
		log.Printf("ERROR: failed to create actor record in repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := ToActorResponse(actor)

	return res, nil
}

func (s *Service) GetActor(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error) {
	const op = "actor.Service.GetActors"

	token := ctx.Value(util.JWTContextKey("jwt")).(string)
	uc, err := util.ParseUserClaims(token, s.signingKey)
	if err != nil {
		log.Printf("ERROR: bad token\n")

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Printf("ERROR: jwt signature is invalid\n")
		}
		if errors.Is(err, util.ErrUnknownClaimsType) {
			log.Printf("ERROR: unknown token claims\n")
		}

		return nil, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	if !uc.IsAdmin {
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

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

	token := ctx.Value(util.JWTContextKey("jwt")).(string)
	uc, err := util.ParseUserClaims(token, s.signingKey)
	if err != nil {
		log.Printf("ERROR: bad token\n")

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Printf("ERROR: jwt signature is invalid\n")
		}
		if errors.Is(err, util.ErrUnknownClaimsType) {
			log.Printf("ERROR: unknown token claims\n")
		}

		return nil, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	if !uc.IsAdmin {
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	id, err := strconv.ParseInt(req.ID, 10, 32)
	if err != nil {
		log.Printf("ERROR: failed id parameter conversion (string -> int32)\n")
		return nil, fmt.Errorf("%s: %w", op, ErrIdInvalid)
	}

	err = ValidateDate(req.Info.Birthday, &util.ValidationError{})
	if err != nil {
		if !errors.Is(err, ErrDateEmpty) {
			log.Printf("ERROR: failed request validation\n")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	actor := ToActor(req.Info)

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

	token := ctx.Value(util.JWTContextKey("jwt")).(string)
	uc, err := util.ParseUserClaims(token, s.signingKey)
	if err != nil {
		log.Printf("ERROR: bad token\n")

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			log.Printf("ERROR: jwt signature is invalid\n")
		}
		if errors.Is(err, util.ErrUnknownClaimsType) {
			log.Printf("ERROR: unknown token claims\n")
		}

		return nil, fmt.Errorf("%s: %w", op, ErrTokenInvalid)
	}

	if !uc.IsAdmin {
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

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
