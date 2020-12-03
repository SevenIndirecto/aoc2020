package main

import (
	"testing"
)

type Fixtures struct {
	Dx       int
	Dy       int
	Expected int
}

func TestTreesOnSlope(t *testing.T) {
	treeMap := LoadMap("aoc03_ex1.txt")
	fixtures := []Fixtures{
		{1, 1, 2},
		{3, 1, 7},
		{5, 1, 3},
		{7, 1, 4},
		{1, 2, 2},
	}

	for _, fixture := range fixtures {
		got := TreesOnSlope(fixture.Dx, fixture.Dy, treeMap)
		if got != fixture.Expected {
			t.Errorf("Test (%v) got %d expected %d", fixture, got, fixture.Expected)
		}
	}
}
