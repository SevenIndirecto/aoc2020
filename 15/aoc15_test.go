package main

import (
	"testing"
)

type Fixture struct {
	Turn int
	Expected int
}

func TestGetNthByStep(t *testing.T) {
	fixtures := []Fixture{
		{4, 0}, {5, 3}, {6, 3},
		{7, 1}, {8, 0}, {9, 4},
		{10, 0}, {2020, 436},
	}

	m := NewMemorizer([]int{0, 3, 6})
	for _, f := range fixtures {
		got := m.GetNthNumberSpoken(f.Turn)

		if got != f.Expected {
			t.Errorf("On turn %d got %d expected %d", f.Turn, got, f.Expected)
		}
	}
}

type FixtureNth struct {
	Numbers []int
	Turn int
	Expected int
}

func TestGetNth(t *testing.T) {
	fixtures := []FixtureNth{
		{[]int{1,3,2}, 2020, 1},
		{[]int{2,1,3}, 2020, 10},
		{[]int{1,2,3}, 2020, 27},
		{[]int{2,3,1}, 2020, 78},
		{[]int{3,2,1}, 2020, 438},
		{[]int{3,1,2}, 2020, 1836},
		{[]int{0,3,6}, 3e7, 175594},
		{[]int{3,1,2}, 3e7, 362},
	}

	for _, f := range fixtures {
		m := NewMemorizer(f.Numbers)
		got := m.GetNthNumberSpoken(f.Turn)

		if got != f.Expected {
			t.Errorf("[%dth] For %v got %d expected %d", f.Turn, f.Numbers, got, f.Expected)
		}
	}
}
