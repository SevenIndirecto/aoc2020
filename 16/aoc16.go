package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type FieldRule struct {
	Name string
	Min1 int
	Max1 int
	Min2 int
	Max2 int
}

func (rule *FieldRule) Validates(value int) bool {
	return value >= rule.Min1 && value <= rule.Max1 || value >= rule.Min2 && value <= rule.Max2
}

type Ticket struct {
	Values              []int
	InvalidValueIndexes []int
}

type Notes struct {
	Rules                []FieldRule
	MyTicket             Ticket
	NearbyTickets        []Ticket
	ValidTickets         []Ticket
	FieldCandidates      []map[string]bool
	SolvedFieldRuleNames map[int]string
}

func (notes *Notes) ValidateNearbyTickets() {
	for i := range notes.NearbyTickets {
		notes.Validate(&notes.NearbyTickets[i])
	}
}

func (notes *Notes) GetScanningRateError() int {
	sum := 0
	for _, ticket := range notes.NearbyTickets {
		for _, invalidFieldIndex := range ticket.InvalidValueIndexes {
			sum += ticket.Values[invalidFieldIndex]
		}
	}
	return sum
}

func (notes *Notes) GetTicketSignature(ticket Ticket) int {
	// mul fields that start with departure
	mul := 1
	for fieldIndex, ruleName := range notes.SolvedFieldRuleNames {
		found, _ := regexp.MatchString(`^departure`, ruleName)
		if found {
			mul *= ticket.Values[fieldIndex]
			fmt.Printf("Found rule %s with value %d, mul now %d\n", ruleName, ticket.Values[fieldIndex], mul)
		}
	}
	return mul
}

func (notes *Notes) Validate(ticket *Ticket) {
	for i, value := range ticket.Values {
		foundFirstValidRule := false
		for _, rule := range notes.Rules {
			if rule.Validates(value) {
				foundFirstValidRule = true
				break
			}
		}

		if !foundFirstValidRule {
			ticket.InvalidValueIndexes = append(ticket.InvalidValueIndexes, i)
		}
	}
	if len(ticket.InvalidValueIndexes) < 1 {
		notes.ValidTickets = append(notes.ValidTickets, *ticket)
	}
}

func (notes *Notes) BuildCandidateList() {
	for _, rule := range notes.Rules {
		for fieldIndex := range notes.MyTicket.Values {
			// Must apply to each valid ticket
			appliesToEachTicket := true
			for _, ticket := range notes.ValidTickets {
				if !rule.Validates(ticket.Values[fieldIndex]) {
					appliesToEachTicket = false
					break
				}
			}

			if appliesToEachTicket {
				notes.FieldCandidates[fieldIndex][rule.Name] = true
			}
		}
	}
}

func (notes *Notes) MatchRulesToFields() {
	// Init SolvedFields
	notes.SolvedFieldRuleNames = make(map[int]string)

	for ; ; {
		// Find first field with 1 rule candidate
		solvedFieldIndex := -1
		var solvedFieldRuleName string

		for i, candidateRules := range notes.FieldCandidates {
			if len(candidateRules) == 1 {
				solvedFieldIndex = i
				// Only one entry, but extract key name
				for ruleName, _ := range candidateRules {
					solvedFieldRuleName = ruleName
					break
				}
				break
			}
		}

		// Did not find any new candidates, assume solved
		if solvedFieldIndex == -1 {
			// done
			break
		}

		// Remove rule as candidate from all
		for i := range notes.FieldCandidates {
			delete(notes.FieldCandidates[i], solvedFieldRuleName)
		}
		notes.SolvedFieldRuleNames[solvedFieldIndex] = solvedFieldRuleName
	}

	// Quick sanity check
	for i, candidates := range notes.FieldCandidates {
		if len(candidates) > 0 {
			panic(fmt.Sprintf("Found field [%d] with unsolved candidates %v", i, candidates))
		}
	}
}

func NewNotes(noteStr string) Notes {
	// Parse rules
	lines := strings.Split(noteStr, "\n")
	notes := Notes{}

	i := 0
	// Parse rules
	var rules []FieldRule
	for ; i < len(lines); i++ {
		line := lines[i]
		if len(line) < 2 {
			break
		}

		s := strings.Split(line, ": ")
		s2 := strings.Split(s[1], " or ")
		s3 := strings.Split(s2[0], "-")
		min1, _ := strconv.Atoi(s3[0])
		max1, _ := strconv.Atoi(s3[1])
		s4 := strings.Split(s2[1], "-")
		min2, _ := strconv.Atoi(s4[0])
		max2, _ := strconv.Atoi(s4[1])

		rule := FieldRule{
			Name: s[0],
			Min1: min1,
			Max1: max1,
			Min2: min2,
			Max2: max2,
		}

		rules = append(rules, rule)
	}
	notes.Rules = rules

	i++
	if lines[i] != "your ticket:" {
		panic("Failed to parse notes")
	}
	i++

	parseTicket := func(ticketStr string) Ticket {
		fields := strings.Split(ticketStr, ",")
		var ticket []int
		for _, f := range fields {
			val, _ := strconv.Atoi(f)
			ticket = append(ticket, val)
		}
		return Ticket{Values: ticket}
	}
	notes.MyTicket = parseTicket(lines[i])

	i += 3
	for ; i < len(lines); i++ {
		line := lines[i]
		if len(line) < 2 {
			break
		}
		notes.NearbyTickets = append(notes.NearbyTickets, parseTicket(line))
	}

	// init field candidates
	notes.FieldCandidates = []map[string]bool{}
	for range notes.MyTicket.Values {
		notes.FieldCandidates = append(notes.FieldCandidates, make(map[string]bool))
	}

	return notes
}

func main() {
	dat, err := ioutil.ReadFile("aoc16.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	notes := NewNotes(txt)
	notes.ValidateNearbyTickets()
	fmt.Println("Part one:", notes.GetScanningRateError())

	notes.BuildCandidateList()
	notes.MatchRulesToFields()
	fmt.Println("Part two:", notes.GetTicketSignature(notes.MyTicket))
}
