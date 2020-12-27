package main

import (
	"testing"
)

type Fixture struct {
	Moves int
	Expected string
}

func TestPartOne(t *testing.T) {
	fixtures := []Fixture{
		{10, "92658374"},
		{100, "67384529"},
	}

	for _, f := range fixtures {
		cg := NewCupGame("389125467", 9)
		cg.Play(f.Moves)
		got := cg.GetPartOneSig()

		if got != f.Expected {
			t.Errorf("Got %s expected %s", got, f.Expected)
		}
	}
}

func TestPartTwo(t *testing.T) {
	expected := 149245887792
	cg := NewCupGame("389125467", 1000000)
	cg.Play(10000000)

	got := cg.GetPartTwoSig()

	if got != expected {
		t.Errorf("Got %d expected %d", got, expected)
	}
}
