package main

import (
	"testing"
)

type Fixtures struct {
	Mass     int
	Expected int
}

func TestFuel(t *testing.T) {
	fixtures := []Fixtures{
		{12, 2},
	}
	for _, fixture := range fixtures {
		got := Fuel(fixture.Mass)
		if got != fixture.Expected {
			t.Errorf("Fuel(%d) = %d; want %d", fixture.Mass, got, fixture.Expected)
		}
	}
}
