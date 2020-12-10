package main

import (
	"testing"
)

type Fixture struct {
	Path string
	Expected int
}

func TestAdapterBag_MapJoltDifference(t *testing.T) {
	fixtures := []Fixture{
		{"aoc10_test1.txt",  7 * 5},
		{"aoc10_test2.txt",  22 * 10},
	}

	for _, fixture := range fixtures {
		bag := AdapterBag{}
		bag.Init(fixture.Path)
		got := bag.MapJoltDifference()

		if got != fixture.Expected {
			t.Errorf("Checking[%s], got: %d expected: %d", fixture.Path, got, fixture.Expected)
		}
	}
}

func TestAdapterBag_GetOptionCount(t *testing.T) {
	fixtures := []Fixture{
		{"aoc10_test1.txt",  8},
		{"aoc10_test2.txt",  19208},
	}

	for _, fixture := range fixtures {
		bag := AdapterBag{}
		bag.Init(fixture.Path)
		got := bag.GetOptionCount(0)

		if got != fixture.Expected {
			t.Errorf("Checking[%s], got: %d expected: %d", fixture.Path, got, fixture.Expected)
		}
	}
}
