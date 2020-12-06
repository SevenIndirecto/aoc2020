package main

import (
	"testing"
)

type Fixtures struct {
	Group    Group
	Expected int
}

func TestGetYesToAnyCount(t *testing.T) {
	fixtures := []Fixtures{
		{Group{"abc"}, 3},
		{Group{"a", "b", "c"}, 3},
		{Group{"ab", "ac"}, 3},
		{Group{"a", "a", "a", "a"}, 1},
		{Group{"b"}, 1},
	}

	for _, fixture := range fixtures {
		got := GetYesToAnyCount(fixture.Group)

		if got != fixture.Expected {
			t.Errorf("GetYesToAnyCount(%v) = %d; want %d", fixture.Group, got, fixture.Expected)
		}
	}
}

func TestGetYesToAllCount(t *testing.T) {
	fixtures := []Fixtures{
		{Group{"abc"}, 3},
		{Group{"a", "b", "c"}, 0},
		{Group{"ab", "ac"}, 1},
		{Group{"a", "a", "a", "a"}, 1},
		{Group{"b"}, 1},
		{Group{"kend", "endk"}, 4},
	}

	for _, fixture := range fixtures {
		got := GetYesToAllCount(fixture.Group)

		if got != fixture.Expected {
			t.Errorf("GetYesToAllCount(%v) = %d; want %d", fixture.Group, got, fixture.Expected)
		}
	}
}
