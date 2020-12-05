package main

import (
	"reflect"
	"testing"
)

type Fixtures struct {
	Pass     string
	Expected Seat
	SeatId int
}

func TestToSeat(t *testing.T) {
	fixtures := []Fixtures{
		{"BFFFBBFRRR", Seat{70, 7}, 567},
		{"FFFBBBFRRR", Seat{14, 7}, 119},
		{"BBFFBBFRLL", Seat{102, 4}, 820},
	}
	for _, fixture := range fixtures {
		got := ToSeat(fixture.Pass)
		if !reflect.DeepEqual(got, fixture.Expected) || got.SeatId() != fixture.SeatId {
			t.Errorf("ToSeat(%s) = %v; want %v", fixture.Pass, got, fixture.Expected)
		}
	}
}
