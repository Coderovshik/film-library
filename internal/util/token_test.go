package util

import "testing"

const (
	secret = "mock_secret"
)

func TestParseUserClaims(t *testing.T) {
	uc := NewUserClaims(1, "user", false)
	ss, err := NewJWTSignedString(uc, secret)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	parsedUC, err := ParseUserClaims(ss, secret)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	if !(parsedUC.ID == uc.ID && parsedUC.Username == uc.Username && parsedUC.IsAdmin == uc.IsAdmin) {
		t.Errorf("Expected %+v, got %+v", uc, parsedUC)
	}
}
