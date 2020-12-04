package main

import (
	"testing"
)

type Fixture struct {
	Path string
	Expected int
}

func TestValidate(t *testing.T) {
	fixtures := []Fixture{
		{"aoc4_test_invalid.txt", 0},
		{"aoc4_test_valid.txt", 4},
	}
	for _, fixture := range fixtures {
		passports := LoadPassports(fixture.Path)
		_, valid := Validate(passports)

		if valid != fixture.Expected {
			t.Errorf(
				"[%s] Got %d valid passports, expected %d",
				fixture.Path,
				valid,
				fixture.Expected,
			)
		}
	}
}
