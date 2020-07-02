package main

import "testing"

func TestRangeString(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Mod      int
		Expected []int
	}{
		{
			"all hours",
			"0-23",
			24,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		},
		{
			"first hour lower than second",
			"4-20",
			24,
			[]int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		{
			"first hour later than second (overnight)",
			"22-4",
			24,
			[]int{22, 23, 0, 1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		r, err := rangeString(test.Input, test.Mod)
		if err != nil {
			t.Errorf("failed test '%s': %v", test.Name, err)
		}
		if !areIntSlicesEqual(r, test.Expected) {
			t.Errorf("results didn't match for '%s'. Expected '%v', got '%v'.", test.Name, test.Expected, r)
		}
	}
}

func TestRangeMonths(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected []int
	}{
		{
			"all months",
			"January-December",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
		{
			"march through september",
			"mar-Sept",
			[]int{2, 3, 4, 5, 6, 7, 8},
		},
		{
			"december through february",
			"dec-February",
			[]int{11, 0, 1},
		},
	}

	for _, test := range tests {
		r, err := rangeMonths(test.Input)
		if err != nil {
			t.Errorf("failed test '%s': %v", test.Name, err)
		}
		if !areIntSlicesEqual(r, test.Expected) {
			t.Errorf("results didn't match for '%s'. Expected '%v', got '%v'.", test.Name, test.Expected, r)
		}
	}
}

func areIntSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
