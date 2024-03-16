package util

import "testing"

func TestQueryableObject_IsEmpty(t *testing.T) {
	qo := &QueryableObject{}

	res := qo.IsEmpty()
	if !res {
		t.Errorf("Expected %t, got %t", true, res)
	}

	qo.Add("obj", 1)
	res = qo.IsEmpty()
	if res {
		t.Errorf("Expected %t, got %t", false, res)
	}
}

func TestQueryableObject_Args(t *testing.T) {
	qo := &QueryableObject{}
	qo.Add("arg1", 1)
	qo.Add("arg2", 2)

	expRes := "arg1 = $1, arg2 = $2"
	res := qo.Args(1)
	if res != expRes {
		t.Errorf("Expected %s, got %s", expRes, res)
	}
}
