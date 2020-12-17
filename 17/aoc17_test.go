package main

import (
	"reflect"
	"testing"
)

const INITIAL_STATE = `.#.
..#
###
`

func TestNewPocketDimension(t *testing.T) {
	got := NewPocketDimension(INITIAL_STATE)
	expected := PocketDimension{
		Cycle: 0,
		Cubes: map[Point]bool{
			Point{0, 0, 0, 0}: false,
			Point{0, 1, 0, 0}: false,
			Point{0, 2,0, 0}: true,
			Point{1, 0, 0, 0}: true,
			Point{1, 1, 0, 0}: false,
			Point{1, 2, 0, 0}: true,
			Point{2, 0, 0, 0}: false,
			Point{2, 1, 0, 0}: true,
			Point{2, 2, 0, 0}: true,
			Point{}: false,
			Point{}: false,
			Point{}: false,
			Point{}: false,
			Point{}: false,
		},
		ConstraintX: Constraint{0, 2},
		ConstraintY: Constraint{0, 2},
		ConstraintZ: Constraint{0, 0},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Did not manage to load initial state got %v, expected %v", got, expected)
	}
}

type FixtureNeighborState struct {
	Point Point
	Expected map[bool]int
}

func TestPocketDimension_GetNeighborState(t *testing.T) {
	pd := NewPocketDimension(INITIAL_STATE)
	pd.CreateSnapshot()

	fixtures := []FixtureNeighborState{
		{Point{0, 0,0, 0}, map[bool]int{ACTIVE: 1, INACTIVE: 79}},
		{Point{2, 2,0, 0}, map[bool]int{ACTIVE: 2, INACTIVE: 78}},
		{Point{2, 2,1, 0}, map[bool]int{ACTIVE: 3, INACTIVE: 77}},
	}

	for _, f := range fixtures {
		got := pd.GetNeighborState(f.Point)

		if !reflect.DeepEqual(got, f.Expected) {
			t.Errorf("Invalid neighbor state got %v, expected %v", got, f.Expected)
		}
	}
}

type FixtureNewState struct {
	IsCurrentlyActive bool
	NeighborState map[bool]int
	Expected bool
}

func TestGetNewState(t *testing.T) {
	fixtures := []FixtureNewState{
		{ACTIVE, map[bool]int{ACTIVE: 3, INACTIVE: 0}, ACTIVE},
		{ACTIVE, map[bool]int{ACTIVE: 2, INACTIVE: 10}, ACTIVE},
		{ACTIVE, map[bool]int{ACTIVE: 1, INACTIVE: 5}, INACTIVE},
		{INACTIVE, map[bool]int{ACTIVE: 3, INACTIVE: 0}, true},
		{INACTIVE, map[bool]int{ACTIVE: 4, INACTIVE: 4}, false},
	}

	for _, f := range fixtures {
		got := GetNewState(f.IsCurrentlyActive, f.NeighborState)

		if got != f.Expected {
			t.Errorf("Invalid new state got %v expected %v", got, f.Expected)
		}
	}
}

type Fixture struct {
	Mode int
	Expected int
}

func TestCycleRuns(t *testing.T) {
	fixtures := []Fixture{
		{MODE_PART_ONE, 112},
		{MODE_PART_TWO, 848},
	}

	for _, f := range fixtures {
		pd := NewPocketDimension(INITIAL_STATE)
		pd.Mode = f.Mode

		for pd.Cycle < 6 {
			pd.ExecuteCycle()
		}

		got := pd.GetActiveCubeCount()

		if got != f.Expected {
			t.Errorf("Test Part One got %d, expected %d", got, f.Expected)
		}
	}
}

