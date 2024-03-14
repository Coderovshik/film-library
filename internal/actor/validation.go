package actor

import (
	"errors"
	"time"

	"github.com/Coderovshik/film-library/internal/util"
)

var (
	ErrDateEmpty = errors.New("date is empty")
)

func ValidateActorInfo(ai *ActorInfo) error {
	ve := &util.ValidationError{}

	if len(ai.Name) == 0 {
		ve.AddViolation("name of length 0")
	}

	if len(ai.Sex) == 0 {
		ve.AddViolation("sex of length 0")
	}

	ValidateDate(ai.Birthday, ve)

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateDate(date string, ve *util.ValidationError) error {
	if len(date) == 0 {
		return ErrDateEmpty
	}

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		ve.AddViolation("incorrect date format (expected format: 2006-01-02)")
		return ve
	}

	return nil
}
