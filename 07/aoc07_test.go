package main

import (
	"reflect"
	"testing"
)

type ParseFixture struct {
	Rule     string
	Expected Bag
}

type RequiredFixture struct {
	Path     string
	Expected int
}

func TestCanContainColor(t *testing.T) {
	fixturePath := "aoc07_test1.txt"
	expected := 4

	bags := LoadRules(fixturePath)
	got := 0
	for _, bag := range bags {
		if CanContainColor(bag, "shiny gold", bags) {
			got++
		}
	}
	if got != expected {
		t.Errorf("TestFixture[%s], got count %d, want %d", fixturePath, got, expected)
	}
}

func TestGetBagsRequiredForBag(t *testing.T) {
	fixtures := []RequiredFixture{
		{"aoc07_test1.txt", 32},
		{"aoc07_test2.txt", 126},
	}

	for _, fixture := range fixtures {
		bags := LoadRules(fixture.Path)
		shinyBag := bags["shiny gold"]
		got := GetBagsRequiredForBag(shinyBag, bags)

		if got != fixture.Expected {
			t.Errorf("TestFixture[%s], got count %d, want %d",
				fixture.Path, got, fixture.Expected,
			)
		}
	}
}

func TestParseBagRule(t *testing.T) {
	fixtures := []ParseFixture{
		{
			Rule: "light red bags contain 1 bright white bag, 2 muted yellow bags.",
			Expected: Bag{
				Color: "light red",
				Capacity: []BagCapacity{
					{Qty: 1, Color: "bright white"},
					{Qty: 2, Color: "muted yellow"},
				},
				KnownToContainTarget: false,
			},
		},
		{
			Rule: "faded blue bags contain no other bags.",
			Expected: Bag{
				Color:                "faded blue",
				KnownToContainTarget: false,
			},
		},
	}

	for _, fixture := range fixtures {
		got := ParseBagRule(fixture.Rule)

		if !reflect.DeepEqual(got, fixture.Expected) {
			t.Errorf("ParseBagRule(%s) = %v, want %v", fixture.Rule, got, fixture.Expected)
		}
	}
}
