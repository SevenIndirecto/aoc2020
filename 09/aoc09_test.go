package main

import (
	"reflect"
	"testing"
)

type FixtureIncrement struct {
	Numbers []int
	PreambleSize int
	Increments int
	ExpectedSums PrecomputedSums
}

func TestXMAS_IncrementPos(t *testing.T) {
	fixtures := []FixtureIncrement{
		{[]int{1, 2, 3, 4, 5, 6}, 3, 2, PrecomputedSums{
			7: 1, 8: 1, 9: 1,
		}},
		{[]int{1, 2, 3, 5, 5, 6}, 3, 2, PrecomputedSums{
			8: 2, 10: 1,
		}},
	}

	for _, fixture := range fixtures {
		xmas := XMAS{
			Numbers: fixture.Numbers,
			PreambleSize: fixture.PreambleSize,
		}
		xmas.Init()
		for i := 0; i < fixture.Increments; i++ {
			xmas.IncrementPos()
		}
		got := xmas.PrecomputedSums

		if !reflect.DeepEqual(got, fixture.ExpectedSums) {
			t.Errorf("%v got %v; want %v", fixture, got, fixture.ExpectedSums)
		}
	}
}

type FixtureValidPos struct {
	NumberToCheck int
	Expected bool
}

func TestXMAS_IsCurrentPosValid(t *testing.T) {
	fixtures := []FixtureValidPos{
		{26, true},
		{49, true},
		{100, false},
		{50, false},
	}

	numbers := [26]int{}
	for i := 0; i < 25; i++ {
		numbers[i] = i+1
	}

	for _, fixture := range fixtures {
		numbers[25] = fixture.NumberToCheck

		xmas := XMAS{
			Numbers: numbers[:],
			PreambleSize: 25,
		}
		xmas.Init()
		got := xmas.IsCurrentPosValid()
		if got != fixture.Expected {
			t.Errorf("Checking %d, got %v expected %v", fixture.NumberToCheck, got, fixture.Expected)
		}
	}
}

type Fixture struct {
	Path string
	PreambleSize int
	Expected int
}

func TestFindFirstInvalid(t *testing.T) {
	fixtures := []Fixture{
		{"aoc09_test1.txt", 5, 127},
	}

	for _, fixture := range fixtures {
		got, _ := FindFirstInvalid(fixture.Path, fixture.PreambleSize)

		if got != fixture.Expected {
			t.Errorf("Checking[%s], got: %d expected: %d", fixture.Path, got, fixture.Expected)
		}
	}
}

func TestFindContiguous(t *testing.T) {
	fixtures := []Fixture{
		{"aoc09_test1.txt", 5, 62},
	}

	for _, fixture := range fixtures {
		target, xmas := FindFirstInvalid(fixture.Path, fixture.PreambleSize)
		got := FindContiguous(xmas.Numbers, target)

		if got != fixture.Expected {
			t.Errorf("Checking[%s], got: %d expected: %d", fixture.Path, got, fixture.Expected)
		}
	}
}
