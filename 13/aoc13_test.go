package main

import (
	"testing"
)

func TestGetFirstBus(t *testing.T) {
	start := 939
	buses := []int{7, 13, 59, 31, 19}
	expected := 295
	bus, wait, _ := GetFirstBus(buses, start)
	got := wait * bus

	if got != expected {
		t.Errorf("Got %d expected %d", got, expected)
	}
}

type FixtureFirstAndPeriod struct {
	Period int
	A int
	B int
	Offset int
	ExpectedFirst int
	ExpectedPeriod int
}

func TestFindFirstAndPeriod(t *testing.T) {
	fixtures := []FixtureFirstAndPeriod{
		{1, 2, 3, 1, 2, 6},
		{1, 3, 5, 1, 9, 15},
		{1, 3, 5, 2, 3, 15},
		{1, 17, 13, 2, 102, 221},
		{221, 102, 19, 3, 3417, 4199},
	}

	for _, f := range fixtures {
		gotFirst, gotPeriod := FindFirstAndPeriod(f.Period, f.A, f.B, f.Offset)

		if gotFirst != f.ExpectedFirst || gotPeriod != f.ExpectedPeriod {
			t.Errorf(
				"FindFirstAndPeriod(%v) got %d, %d expected %d, %d",
				f,
				gotFirst,
				gotPeriod,
				f.ExpectedFirst,
				f.ExpectedPeriod,
			)
		}
	}
}

type FixtureSolve struct {
	Buses []Bus
	Expected int
}

func TestSolveSchedule(t *testing.T) {
	fixture := []FixtureSolve{
		{
			Buses: []Bus{{17, 0}, {13, 2}, {19, 3}},
			Expected: 3417,
		},
		{
			Buses: []Bus{{67, 0}, {7, 1}, {59, 2}, {61, 3}},
			Expected: 754018,
		},
		{
			Buses: []Bus{{67, 0}, {7, 2}, {59, 3}, {61, 4}},
			Expected: 779210,
		},
		{
			Buses: []Bus{{67, 0}, {7, 1}, {59, 3}, {61, 4}},
			Expected: 1261476,
		},
		{
			Buses: []Bus{{1789, 0}, {37, 1}, {47, 2}, {1889, 3}},
			Expected: 1202161486,
		},
	}

	for _, f := range fixture {
		got := SolveSchedule(f.Buses)

		if got != f.Expected {
			t.Errorf("Solve(%v) got %d expected %d", f.Buses, got, f.Expected)
		}
	}
}
