package main

import "testing"

func TestRng(t *testing.T) {
	tests := []struct {
		Name     string
		Min      int
		Max      int
		Expected []int
	}{
		{
			"0-11",
			0,
			11,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
		{
			"11-0",
			11,
			0,
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
		{
			"3 only",
			3,
			3,
			[]int{3},
		},
		{
			"3-21",
			3,
			21,
			[]int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
		},
	}

	for _, test := range tests {
		got := rng(test.Min, test.Max)
		if !areIntSlicesEqual(got, test.Expected) {
			t.Errorf("failed test '%s': expected '%v', got '%v'", test.Name, test.Expected, got)
		}
	}
}

func TestParseMonths(t *testing.T) {
	tests := []struct {
		Input       string
		Expected    []int
		ShouldError bool
	}{
		{
			"All",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			false,
		},
		{
			"all",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			false,
		},
		{
			"June",
			[]int{5},
			false,
		},
		{
			"May, June, July, August, September, October",
			[]int{4, 5, 6, 7, 8, 9},
			false,
		},
		{
			"January, February, March, November, December",
			[]int{10, 11, 0, 1, 2},
			false,
		},
		{
			"January, February, March, April, July, August, September, November, December",
			[]int{0, 1, 2, 3, 6, 7, 8, 10, 11},
			false,
		},
		{
			"January, February, March, April, May, December",
			[]int{0, 1, 2, 3, 4, 11},
			false,
		},
		{
			"All except July, August",
			[]int{0, 1, 2, 3, 4, 5, 8, 9, 10, 11},
			false,
		},
		{
			"all except september, october, november",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 11},
			false,
		},
		{
			"all except jan, jun, september",
			[]int{1, 2, 3, 4, 6, 7, 9, 10, 11},
			false,
		},
		{
			"fhrblig",
			[]int{},
			true,
		},
		{
			"All except jul",
			[]int{0, 1, 2, 3, 4, 5, 7, 8, 9, 10, 11},
			false,
		},
	}

	for _, test := range tests {
		got, err := parseMonths(test.Input)
		if err != nil && !test.ShouldError {
			t.Errorf("failed test '%s', shouldn't have errored but got: %v", test.Input, err)
		} else {
			if !areIntSlicesEqual(got, test.Expected) {
				t.Errorf("failed test '%s': expected '%v', got '%v'", test.Input, test.Expected, got)
			}
		}
	}
}

func TestInvertMonths(t *testing.T) {
	tests := []struct {
		Input    []string
		Expected []string
	}{
		{
			[]string{"Jan", "February"},
			[]string{"mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"},
		},
		{
			[]string{"November", "December", "Jan", "feb"},
			[]string{"mar", "apr", "jun", "may", "jul", "aug", "sep", "oct"},
		},
	}

	for _, test := range tests {
		got := invertMonths(test.Input)
		if !areStringSlicesEqual(got, test.Expected) {
			t.Errorf("failed test '%v': expected '%v', got '%v'", test.Input, test.Expected, got)
		}
	}
}

func TestParseHours(t *testing.T) {
	tests := []struct {
		Input       string
		Expected    []int
		ShouldError bool
	}{
		{
			"7AM-9AM",
			[]int{7, 8},
			false,
		},
		{
			"7AM-4PM",
			[]int{7, 8, 9, 10, 11, 12, 13, 14, 15},
			false,
		},
		{
			"All",
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
			false,
		},
		{
			"10AM-12PM",
			[]int{10, 11},
			false,
		},
		{
			"12PM-4PM",
			[]int{12, 13, 14, 15},
			false,
		},
		{
			"12AM-4AM",
			[]int{0, 1, 2, 3},
			false,
		},
		{
			"9PM-3AM",
			[]int{21, 22, 23, 0, 1, 2},
			false,
		},
		{
			"10PM-2AM, 8AM-10AM",
			[]int{22, 23, 0, 1, 8, 9},
			false,
		},
	}

	for _, test := range tests {
		got, err := parseHours(test.Input)
		if test.ShouldError && err == nil {
			t.Errorf("Failed test '%s' - should have errored but didn't", test.Input)
			continue
		}
		if err != nil {
			t.Errorf("Failed test '%s', should not have errored but got '%v'", test.Input, err)
			continue
		}
		if !areIntSlicesEqual(got, test.Expected) {
			t.Errorf("Failed test '%s': Expected '%v', got '%v'", test.Input, test.Expected, got)
		}
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		Input       string
		Expected    int
		ShouldError bool
	}{
		{
			"7AM",
			7,
			false,
		},
		{
			"12AM",
			0,
			false,
		},
		{
			"4PM",
			16,
			false,
		},
		{
			"7",
			7,
			true,
		},
		{
			"12PM",
			12,
			false,
		},
	}

	for _, test := range tests {
		got, err := parseTime(test.Input)
		if test.ShouldError && err == nil {
			t.Errorf("Failed test '%s', should have errored but didn't", test.Input)
		}
		if !test.ShouldError {
			if err != nil {
				t.Errorf("Failed test '%s', shouldn't have errored but got '%v'", test.Input, err)
			} else {
				if got != test.Expected {
					t.Errorf("Failed test '%s': expected %d, got %d", test.Input, test.Expected, got)
				}
			}
		}
	}
}

func areStringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	itemsInA := make(map[string]bool)
	for i := range a {
		itemsInA[a[i]] = true
	}

	for i := range b {
		if !itemsInA[b[i]] {
			return false
		}
	}

	return true
}

func areIntSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	itemsInA := make(map[int]bool)
	for i := range a {
		itemsInA[a[i]] = true
	}

	for i := range b {
		if !itemsInA[b[i]] {
			return false
		}
	}

	return true
}
