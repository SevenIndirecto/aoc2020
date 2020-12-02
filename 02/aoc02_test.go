package main

import (
	"reflect"
	"testing"
)

type Fixtures struct {
	Line     string
	Expected bool
}

func TestParseLine(t *testing.T) {
	fixtures := map[string]DbEntry{
		"1-3 a: abcde": {Policy{1, 3, []rune(`a`)[0]}, "abcde"},
	}

	for line, expected := range fixtures {
		got, _ := ParseLine(line)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("ParseLine(%v) = %v expected %v", line, got, expected)
		}
	}
}

func TestIsValidPartTwo(t *testing.T) {
	fixtures := []Fixtures{
		{"1-3 a: abcde", true},
		{"1-3 b: cdefg", false},
		{"2-9 c: ccccccccc", false},
	}

	for _, fixture := range fixtures {
		entry, _ := ParseLine(fixture.Line)
		got := IsValidPartTwo(entry)
		if got != fixture.Expected {
			t.Errorf("IsValidPartTwo(%v) = %v expected %v", fixture.Line, got, fixture.Expected)
		}
	}
}

func TestIsValidPartOne(t *testing.T) {
	fixtures := []Fixtures{
		{"1-3 a: abcde", true},
		{"1-3 b: cdefg", false},
		{"2-9 c: ccccccccc", true},
	}

	for _, fixture := range fixtures {
		entry, _ := ParseLine(fixture.Line)
		got := IsValidPartOne(entry)
		if got != fixture.Expected {
			t.Errorf("IsValidPartOne(%v) = %v expected %v", fixture.Line, got, fixture.Expected)
		}
	}
}
