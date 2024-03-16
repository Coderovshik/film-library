package actor

import (
	"testing"
	"time"
)

func TestToQueryableObject(t *testing.T) {
	bd, _ := time.Parse(time.DateOnly, "1995-05-03")
	a := &Actor{
		Name:     "actor",
		Sex:      "male",
		Birthday: bd,
	}

	qo := ToQueryableObject(a)
	args, values := qo.Args(1), qo.Values()
	if len(values) != 3 || args != "actor_name = $1, sex = $2, birthday = $3" {
		t.Fatalf("Expected %d values with args '%s', got %d values with args '%s'",
			3, "actor_name = $1, sex = $2, birthday = $3", len(values), args)
	}

	a = &Actor{}
	qo = ToQueryableObject(a)
	args, values = qo.Args(1), qo.Values()
	if !qo.IsEmpty() {
		t.Fatalf("Expected empty object, got %d values with args '%s'",
			len(values), args)
	}
}
