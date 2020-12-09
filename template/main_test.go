package main

import (
	"testing"
)

type Fixture struct {
}

func TestFuel(t *testing.T) {
	fixtures := []Fixture{}

	for _, fixture := range fixtures {
		got := Fuel(fixture.Mass)

		if got != fixture.Expected {
			t.Errorf("Fuel(%d) = %d; want %d", fixture.Mass, got, fixture.Expected)
		}
	}
}
