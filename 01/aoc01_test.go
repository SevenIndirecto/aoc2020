package main

import (
	"sort"
	"testing"
)

type Fixtures struct {
	Inputs   []int
	Target   int
	Expected int
}

func TestExpenseReport(t *testing.T) {
	fixtures := []Fixtures{
		{[]int{1721, 979, 366, 299, 675, 1456}, 2020, 514579},
	}
	for _, fixture := range fixtures {
		sort.Sort(sort.Reverse(sort.IntSlice(fixture.Inputs)))
		a, b := ExpenseReport(fixture.Inputs, fixture.Target)
		if a * b != fixture.Expected {
			t.Errorf("ExpenseReport(%d) = %d; want %d", fixture.Inputs, a * b, fixture.Expected)
		}
	}
}

func TestExpenseReportThree(t *testing.T) {
	fixtures := []Fixtures{
		{[]int{1721, 979, 366, 299, 675, 1456}, 2020, 241861950},
	}
	for _, fixture := range fixtures {
		sort.Sort(sort.Reverse(sort.IntSlice(fixture.Inputs)))
		a, b, c := ExpenseReportThree(fixture.Inputs, fixture.Target)
		got := a * b * c
		if got != fixture.Expected {
			t.Errorf("ExpenseReportThree(%d) = %d; want %d", fixture.Inputs, got, fixture.Expected)
		}
	}
}
