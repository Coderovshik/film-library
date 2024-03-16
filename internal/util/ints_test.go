package util

import "testing"

func TestRemoveDuplicateInt(t *testing.T) {
	var tests = []struct {
		name  string
		input []int
		want  []int
	}{
		{
			"No dulicates",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"1 dulicate",
			[]int{1, 2, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"2 dulicates",
			[]int{1, 1, 2, 3, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		},
		{
			"Dulicates only",
			[]int{1, 1, 1, 1, 1},
			[]int{1},
		},
		{
			"Empty",
			[]int{},
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := RemoveDuplicateInt(tt.input)
			if !AreEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

func TestAreEqual(t *testing.T) {
	var tests = []struct {
		name   string
		inputA []int
		inputB []int
		want   bool
	}{
		{
			"Ints equal",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			true,
		},
		{
			"Ints unequal",
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3},
			false,
		},
		{
			"Ints empty",
			[]int{},
			[]int{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := AreEqual(tt.inputA, tt.inputB)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
