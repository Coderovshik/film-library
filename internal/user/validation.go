package user

import "github.com/Coderovshik/film-library/internal/util"

func ValidateCreateUserReuqest(req *CreateUserRequest) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(req.Username) == 0 {
		ve.AddViolation("username of length 0")
	}

	if len(req.Password) == 0 {
		ve.AddViolation("password of length 0")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}

func ValidateLoginReuqest(req *LoginRequest) *util.ValidationError {
	ve := &util.ValidationError{}

	if len(req.Username) == 0 {
		ve.AddViolation("username of length 0")
	}

	if len(req.Password) == 0 {
		ve.AddViolation("password of length 0")
	}

	if ve.NoViolations() {
		return nil
	}

	return ve
}
