package main

import (
	"reflect"
	"testing"
)

const FixtureNotes = `class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12
`

func TestNewNotes(t *testing.T) {
	notes := NewNotes(FixtureNotes)

	if len(notes.NearbyTickets) != 4 || len(notes.MyTicket.Values) != 3 || len(notes.Rules) != 3 {
		t.Errorf("Failed to parse notes, got notes %v", notes)
	}
}

func TestNotes_ValidateNearbyTickets(t *testing.T) {
	notes := NewNotes(FixtureNotes)
	notes.ValidateNearbyTickets()
	expected := 71

	got := notes.GetScanningRateError()

	if got != expected {
		t.Errorf("Failed to validate tickets, got %d expected %d", got, expected)
	}
}

const FixturePartTwo = `class: 0-1 or 4-19
row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9
`

func TestNotes_BuildCandidateList(t *testing.T) {
	expectedCandidates := []map[string]bool{
		{"row": true},
		{"row": true, "class": true},
		{"row": true, "class": true, "seat": true},
	}
	notes := NewNotes(FixturePartTwo)
	notes.ValidateNearbyTickets()
	notes.BuildCandidateList()
	got := notes.FieldCandidates

	if !reflect.DeepEqual(got, expectedCandidates) {
		t.Errorf("Invalid candidates got %v, expected %v", got, expectedCandidates)
	}
}

func TestNotes_MatchRulesToFields(t *testing.T) {
	expectedFieldRules := map[int]string{
		0: "row",
		1: "class",
		2: "seat",
	}
	notes := NewNotes(FixturePartTwo)
	notes.ValidateNearbyTickets()
	notes.BuildCandidateList()
	notes.MatchRulesToFields()
	got := notes.SolvedFieldRuleNames

	if !reflect.DeepEqual(got, expectedFieldRules) {
		t.Errorf("Invalid matches got %v, expected %v", got, expectedFieldRules)
	}
}

func TestNotes_GetTicketSignature(t *testing.T) {
	notesStr := `departure class: 0-1 or 4-19
departure row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9
`

	notes := NewNotes(notesStr)
	notes.ValidateNearbyTickets()
	notes.BuildCandidateList()
	notes.MatchRulesToFields()
	got := notes.GetTicketSignature(notes.MyTicket)
	expected := 132

	if got != expected {
		t.Errorf("Got %d expected %d", got, expected)
	}
}
