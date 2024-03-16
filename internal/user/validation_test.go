package user

import "testing"

func TestValidateCreateUserReuqest(t *testing.T) {
	req := &CreateUserRequest{
		Username: "user",
		Password: "user",
	}

	vErr := ValidateCreateUserReuqest(req)
	if vErr != nil {
		t.Fatalf("No error expected, got %s", vErr.Error())
	}

	req = &CreateUserRequest{
		Username: "",
		Password: "",
	}

	vErr = ValidateCreateUserReuqest(req)
	if vErr.Error() != "username of length 0; password of length 0" {
		t.Fatalf("No %s, got %s", "username of length 0; password of length 0", vErr.Error())
	}
}

func TestValidateLoginReuqest(t *testing.T) {
	req := &LoginRequest{
		Username: "user",
		Password: "user",
	}

	vErr := ValidateLoginReuqest(req)
	if vErr != nil {
		t.Fatalf("No error expected, got %s", vErr.Error())
	}

	req = &LoginRequest{
		Username: "",
		Password: "",
	}

	vErr = ValidateLoginReuqest(req)
	if vErr.Error() != "username of length 0; password of length 0" {
		t.Fatalf("No %s, got %s", "username of length 0; password of length 0", vErr.Error())
	}
}
