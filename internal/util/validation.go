package util

import "strings"

type ValidationError struct {
	violations []string
}

func (ve *ValidationError) NoViolations() bool {
	return len(ve.violations) == 0
}

func (ve *ValidationError) AddViolation(v string) {
	ve.violations = append(ve.violations, v)
}

func (ve *ValidationError) Error() string {
	return strings.Join(ve.violations[:], "; ")
}
