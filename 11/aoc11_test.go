package main

import (
	"reflect"
	"testing"
)

type FixtureValues struct {
	Grid     [][]PointState
	Expected PointState
	Mode     Mode
}

type FixtureNeighbors struct {
	Grid     [][]PointState
	Expected map[PointState]int
}

func TestGetNewValue(t *testing.T) {
	fixtures := []FixtureValues{
		{
			[][]PointState{
				{Occupied, Empty, Occupied},
				{Occupied, Occupied, Occupied},
				{Occupied, Occupied, Occupied},
			},
			Empty,
			PartOneMode,
		},
		{
			[][]PointState{
				{Occupied, Empty, Occupied},
				{Empty, Occupied, Empty},
				{Empty, Empty, Empty},
			},
			Occupied,
			PartOneMode,
		},
		{
			[][]PointState{
				{Empty, Empty, Empty},
				{Empty, Occupied, Empty},
				{Empty, Empty, Empty},
			},
			Occupied,
			PartOneMode,
		},
		{
			[][]PointState{
				{Empty, Empty, Empty},
				{Empty, Empty, Empty},
				{Empty, Empty, Empty},
			},
			Occupied,
			PartOneMode,
		},
		{
			[][]PointState{
				{Empty, Empty},
				{Empty, Empty},
			},
			Occupied,
			PartOneMode,
		},
		{
			[][]PointState{
				{Occupied, Occupied, Empty, Empty, Empty},
				{Empty, Occupied, Floor, Occupied, Empty},
				{Empty, Floor, Floor, Floor, Floor},
				{Empty, Occupied, Empty, Occupied},
			},
			Occupied,
			PartOneMode,
		},
		{
			[][]PointState{
				{Occupied, Occupied, Empty, Empty, Empty},
				{Empty, Occupied, Floor, Occupied, Empty},
				{Empty, Floor, Floor, Floor, Floor},
				{Empty, Occupied, Empty, Occupied},
			},
			Empty,
			PartTwoMode,
		},
	}

	for _, fixture := range fixtures {
		var threshold int
		var adjacentState map[PointState]int

		if fixture.Mode == PartOneMode {
			adjacentState = GetAdjacentState(fixture.Grid, Point{1, 1})
			threshold = 4
		} else {
			adjacentState = GetFirstVisibleChairs(fixture.Grid, Point{1, 1})
			threshold = 5
		}
		got := GetNewValue(adjacentState, fixture.Grid, Point{1, 1}, threshold)

		if got != fixture.Expected {
			t.Errorf("GetNewValue(%v) = %v; want %v", fixture.Grid, got, fixture.Expected)
		}
	}
}

func TestGetAdjacentState(t *testing.T) {
	fixtures := []FixtureNeighbors{
		{
			[][]PointState{
				{Occupied, Empty, Occupied},
				{Occupied, Occupied, Floor},
				{Occupied, Occupied, Occupied},
			},
			map[PointState]int{Occupied: 6, Empty: 1, Floor: 1},
		},
		{
			[][]PointState{
				{Empty, Empty},
				{Floor, Empty},
			},
			map[PointState]int{Occupied: 0, Empty: 2, Floor: 1},
		},
	}

	for _, fixture := range fixtures {
		got := GetAdjacentState(fixture.Grid, Point{1, 1})

		if !reflect.DeepEqual(got, fixture.Expected) {
			t.Errorf("GetAdjacentState(%v) = %v; want %v", fixture.Grid, got, fixture.Expected)
		}
	}
}

type FixtureFirstVisible struct {
	Path     string
	Point    Point
	Expected map[PointState]int
}

func TestGetFirstVisibleChairs(t *testing.T) {
	fixtures := []FixtureFirstVisible{
		{
			"aoc11_test2.txt",
			Point{3, 4},
			map[PointState]int{Occupied: 8, Empty: 0, Floor: 0},
		},
		{
			"aoc11_test3.txt",
			Point{1, 1},
			map[PointState]int{Occupied: 0, Empty: 1, Floor: 0},
		},
		{
			"aoc11_test4.txt",
			Point{3, 3},
			map[PointState]int{Occupied: 0, Empty: 0, Floor: 0},
		},
	}

	for _, fixture := range fixtures {
		sl := SeatLayout{}
		sl.Init(fixture.Path)

		got := GetFirstVisibleChairs(sl.Grid, fixture.Point)

		if !reflect.DeepEqual(got, fixture.Expected) {
			t.Errorf("GetFirstVisibleChairs[%v, %v] = %v; want %v",
				fixture.Path, fixture.Point, got, fixture.Expected,
			)
		}
	}

}

func TestPartOne(t *testing.T) {
	sl := SeatLayout{Mode: PartOneMode}
	sl.Init("aoc11_test1.txt")
	expected := 37

	got := FindEquilibrium(&sl)
	if got != expected {
		t.Errorf("Part one: got %d expected %d", got, expected)
	}
}

func TestPartTwo(t *testing.T) {
	sl := SeatLayout{Mode: PartTwoMode}
	sl.Init("aoc11_test1.txt")
	expected := 26

	got := FindEquilibrium(&sl)
	if got != expected {
		t.Errorf("Part two: got %d expected %d", got, expected)
	}
}
