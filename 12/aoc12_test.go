package main

import (
	"fmt"
	"strings"
	"testing"
)

type FixtureTurn struct {
	InitialDirection int
	Instruction Instruction
	Expected int
}

func TestShip_Turn(t *testing.T) {
	fixtures := []FixtureTurn{
		{North, Instruction{"R", 90}, East},
		{North, Instruction{"L", 90}, West},
		{East, Instruction{"L", 180}, West},
		{East, Instruction{"R", 180}, West},
	}
	ship := Ship{}

	for _, fixture := range fixtures {
		ship.Direction = fixture.InitialDirection
		ship.Turn(fixture.Instruction)
		got := ship.Direction

		if got != fixture.Expected {
			t.Errorf(
				"Turn(init: %d, ins: %v) got %d expected %d",
				fixture.InitialDirection,
				fixture.Instruction,
				got,
				fixture.Expected,
			)
		}
	}
}

type FixtureMoveDir struct {
	InitialPosition Point
	Instruction Instruction
	ExpectedPosition Point
}

func TestShip_MoveInDirection(t *testing.T) {
	fixtures := []FixtureMoveDir{
		{Point{0, 0}, Instruction{"N", 5}, Point{0, 5}},
		{Point{10, -5}, Instruction{"S", 1}, Point{10, -6}},
		{Point{0, 0}, Instruction{"W", 11}, Point{-11, 0}},
		{Point{1, 1}, Instruction{"E", 2}, Point{3, 1}},
	}

	for _, f := range fixtures {
		ship := Ship{Position: f.InitialPosition}
		ship.MoveInDirection(f.Instruction)

		got := ship.Position
		if got != f.ExpectedPosition {
			t.Errorf("MoveInDirection(%v), got %v expected %v", f, got, f.ExpectedPosition)
		}
	}
}

type FixtureMoveForward struct {
	InitialPosition Point
	Direction int
	Distance int
	ExpectedPosition Point
}

func TestShip_MoveForward(t *testing.T) {
	fixtures := []FixtureMoveForward {
		{Point{0, 0}, North, 10, Point{0, 10}},
		{Point{10, -5}, East, 2, Point{12, -5}},
		{Point{0, 0}, South, 20, Point{0, -20}},
		{Point{1, 1}, West, 2, Point{-1, 1}},
	}

	for _, f := range fixtures {
		ship := Ship{Position: f.InitialPosition}
		ship.Direction = f.Direction
		ship.MoveForward(f.Distance)

		got := ship.Position
		if got != f.ExpectedPosition {
			t.Errorf("MoveForward(%v), got %v expected %v", f, got, f.ExpectedPosition)
		}
	}
}

type FixtureRotate struct {
	InitialPosition Point
	Instruction Instruction
	Expected Point
}

func TestShip_RotateWaypoint(t *testing.T) {
	fixtures := []FixtureRotate{
		{Point{1,1}, Instruction{"L", 90}, Point{-1, 1}},
		{Point{1,1}, Instruction{"R", 90}, Point{1, -1}},
		{Point{1,1}, Instruction{"L", 270}, Point{1, -1}},
		{Point{1,1}, Instruction{"R", 180}, Point{-1, -1}},
	}

	for _, f := range fixtures {
		ship := Ship{Waypoint: f.InitialPosition}
		ship.RotateWaypoint(f.Instruction)
		got := ship.Waypoint

		if got != f.Expected {
			t.Errorf("RotateWaypoint(%v) got %v expected %v", f, got, f.Expected)
		}
	}
}

type FixtureMoveToWaypoint struct {
	InitialPosition Point
	Waypoint Point
	Instruction Instruction
	Expected Point
}

func TestShip_MoveToWaypoint(t *testing.T) {
	fixtures := []FixtureMoveToWaypoint{
		{
			Point{0, 0},
			Point{10, 1},
			Instruction{"F", 10},
			Point{100, 10},
		},
		{
			Point{170, 38},
			Point{4, -10},
			Instruction{"F", 11},
			Point{214, -72},
		},
	}

	for _, f := range fixtures {
		ship := Ship{Waypoint: f.Waypoint, Position: f.InitialPosition}
		ship.MoveToWaypoint(f.Instruction)
		got := ship.Position

		if got != f.Expected {
			t.Errorf("MoveToWaypoint(%v) got %v expected %v", f, got, f.Expected)
		}
	}
}

func TestShip_ExecuteInstructions(t *testing.T) {
	linesStr := `F10
N3
F7
R90
F11
`
	lines := strings.Split(linesStr, "\n")

	ship := Ship{Direction: East, Position: Point{0, 0}}
	ship.LoadInstructionSet(lines)
	ship.ExecuteInstructions()

	got := ship.DistanceTravelled()
	expected := 25

	if got != expected {
		fmt.Println(ship)
		t.Errorf("Execute instructions, got %d expected %d", got, expected)
	}
}

func TestShip_ExecuteRealInstructions(t *testing.T) {
	linesStr := `F10
N3
F7
R90
F11
`
	lines := strings.Split(linesStr, "\n")

	ship := Ship{Waypoint: Point{10, 1}, Position: Point{0, 0}}
	ship.LoadInstructionSet(lines)
	ship.ExecuteRealInstructions()

	got := ship.DistanceTravelled()
	expected := 286

	if got != expected {
		fmt.Println(ship)
		t.Errorf("Execute real instructions, got %d expected %d", got, expected)
	}
}
