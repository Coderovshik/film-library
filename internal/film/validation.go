package film

import (
	"regexp"
	"time"

	"github.com/Coderovshik/film-library/internal/util"
)

var validSortQuery = regexp.MustCompile("^(name|rating|releasedate),(asc|desc)$")

func ValidateGetFilmsRequest(req *GetFilmsRequest) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(req.SortQuery) != 0 && !validSortQuery.MatchString(req.SortQuery) {
		ve.AddViolation("incorrect sort query, expect value of pattern: '^(name|rating|releasedate),(asc|desc)$'")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateFormatFilmInfo(fi *FilmInfo, allowNegativeRating bool) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(fi.Name) > 150 {
		ve.AddViolation("name length is more than 150 symbols")
	}

	if len(fi.Description) > 1000 {
		ve.AddViolation("description length is more than 1000 symbols")
	}

	if _, err := time.Parse("2006-01-02", fi.ReleaseDate); err != nil && len(fi.ReleaseDate) != 0 {
		ve.AddViolation("incorrect date format (expected format: 2006-01-02)")
	}

	if (fi.Rating < 0 && !allowNegativeRating) || fi.Rating > 10 {
		ve.AddViolation("incorrect rating, expected: 0 <= rating <= 10")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateEmptyFilmInfo(fi *FilmInfo) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(fi.Name) == 0 {
		ve.AddViolation("name empty")
	}

	if len(fi.ReleaseDate) == 0 {
		ve.AddViolation("date empty (expected format: 2006-01-02)")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}
