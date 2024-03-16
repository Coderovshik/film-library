package util

import "testing"

func TestValidationError_NoViolations(t *testing.T) {
	ve := &ValidationError{}

	res := ve.NoViolations()
	if !res {
		t.Fatalf("Expected %t, got %t", true, res)
	}

	ve.AddViolation("v1")
	ve.AddViolation("v2")

	res = ve.NoViolations()
	if res {
		t.Fatalf("Expected %t, got %t", false, res)
	}
}

func TestValidationError_Error(t *testing.T) {
	ve := &ValidationError{}

	res := ve.Error()
	if res != "" {
		t.Fatalf("Expected empty string, got %s", res)
	}

	ve.AddViolation("v1")
	ve.AddViolation("v2")

	res = ve.Error()
	if res != "v1; v2" {
		t.Fatalf("Expected %s, got %s", "v1; v2", res)
	}
}
