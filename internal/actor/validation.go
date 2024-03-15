package actor

import (
	"time"

	"github.com/Coderovshik/film-library/internal/util"
)

var sexMap = map[string]struct{}{
	"":       struct{}{},
	"female": struct{}{},
	"male":   struct{}{},
}

func ValidateFormatActorInfo(ai *ActorInfo) *util.ValidationError {
	ve := &util.ValidationError{}

	if _, ok := sexMap[ai.Sex]; !ok {
		ve.AddViolation("incorrect sex format (expected one of [male, female])")
	}

	if _, err := time.Parse("2006-01-02", ai.Birthday); err != nil && len(ai.Birthday) != 0 {
		ve.AddViolation("incorrect date format (expected format: 2006-01-02)")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateEmptyActorInfo(ai *ActorInfo) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(ai.Name) == 0 {
		ve.AddViolation("name empty")
	}

	if len(ai.Sex) == 0 {
		ve.AddViolation("sex empty (expected one of [male, female])")
	}

	if len(ai.Birthday) == 0 {
		ve.AddViolation("date empty (expected format: 2006-01-02)")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}
