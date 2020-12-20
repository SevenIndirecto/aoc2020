package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type SubRuleSet []int

type Pattern []uint8

type Rule struct {
	Char     uint8
	SubRules []SubRuleSet
}

func (rule *Rule) isSimple() bool {
	return len(rule.SubRules) < 1
}

const NO_MATCH = -1000

type Matcher struct {
	Rules           map[int]Rule
	ActivePattern   int
	Patterns        []Pattern
	MatchedPatterns []int
	RulePath []int
	Patched bool
	TailRecursionRule int
}

func (m *Matcher) PatternMatchesRule(patternId, ruleId int) bool {
	m.ActivePattern = patternId
	matchedUpTo := m.MatchRule(ruleId, -1)

	return matchedUpTo == len(m.Patterns[patternId])-1
}

func (m *Matcher) MatchRule(ruleId int, matchedUpTo int) int {
	m.RulePath = append(m.RulePath, ruleId)
	rule := m.Rules[ruleId]

	if rule.isSimple() {

		if len(m.Patterns[m.ActivePattern]) <= matchedUpTo+1 {

			m.RulePath = m.RulePath[:len(m.RulePath)-1]
			return NO_MATCH
		}
		if m.Patterns[m.ActivePattern][matchedUpTo+1] == rule.Char {

			m.RulePath = m.RulePath[:len(m.RulePath)-1]
			return matchedUpTo + 1
		} else {

			m.RulePath = m.RulePath[:len(m.RulePath)-1]
			return NO_MATCH
		}
	}

	// For tail recursion patching
	subRulesSets := make([]SubRuleSet, len(rule.SubRules))
	copy(subRulesSets, rule.SubRules)

	for recursiveTries := 0 ; recursiveTries < 15 ; recursiveTries++ {

		for _, subRuleSet := range subRulesSets {
			newMatchedUpTo := matchedUpTo

			// All rules in sub rule set should match in ORDER while moving visited pointer forward, validating string up
			// to that point
			for _, subRuleId := range subRuleSet {
				newMatchedUpTo = m.MatchRule(subRuleId, newMatchedUpTo)

				if newMatchedUpTo == NO_MATCH {
					break
				}
			}

			if ruleId == 0 && newMatchedUpTo != len(m.Patterns[m.ActivePattern]) - 1 {
				break
			}

			if newMatchedUpTo != NO_MATCH {
				// Current sub rules match -> DONE
				return newMatchedUpTo
			}
		}

		// Only special handling for tail recursive subrule 8 of rule 0
		if !m.Patched || ruleId != 0 {
			break // no recursion
		} else {
			subRulesSets[0] = append(SubRuleSet{m.TailRecursionRule}, subRulesSets[0]...)
		}
	}
	m.RulePath = m.RulePath[:len(m.RulePath)-1]
	return NO_MATCH
}

func (m *Matcher) PatchForPartTwo(tailRecursionRule int) {
	// No need to patch rule 8, since we're hacking around it...
	//rule := m.Rules[8]
	//rule.SubRules = append(rule.SubRules, SubRuleSet{42, 8})
	//m.Rules[8] = rule
	m.Patched = true
	m.TailRecursionRule = tailRecursionRule

	rule11 := m.Rules[11]
	rule11.SubRules = append(rule11.SubRules, SubRuleSet{42, 11, 31})
	m.Rules[11] = rule11
}

func NewMatcher(list string) Matcher {
	lines := strings.Split(list, "\n")

	m := Matcher{
		Rules: make(map[int]Rule),
		ActivePattern: -1,
		Patched: false,
	}

	// Parse rules
	patternsStartAt := -1
	for i, line := range lines {
		line := strings.TrimRight(line, " ")
		if len(line) < 2 {
			patternsStartAt = i + 1
			break
		}
		s := strings.Split(line, ": ")
		ruleId, _ := strconv.Atoi(s[0])

		rule := Rule{}
		re := regexp.MustCompile(`^"([a-z])"`)
		match := re.FindStringSubmatch(s[1])
		if match != nil {
			// Found leaf rule
			rule.Char = match[1][0]
		} else {
			subRulesString := strings.Split(s[1], " | ")
			var subRuleSets []SubRuleSet

			for _, subRuleSetString := range subRulesString {
				subRuleSetIds := strings.Split(subRuleSetString, " ")
				var singleSet SubRuleSet
				for _, idStr := range subRuleSetIds {
					id, _ := strconv.Atoi(idStr)
					singleSet = append(singleSet, id)
				}

				subRuleSets = append(subRuleSets, singleSet)
			}
			rule.SubRules = subRuleSets
		}

		m.Rules[ruleId] = rule
	}

	if patternsStartAt < 0 || patternsStartAt >= len(lines) {
		return m
	}

	// Load patterns
	var patterns []Pattern
	for i := patternsStartAt; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			continue
		}
		patterns = append(patterns, Pattern(lines[i]))
	}
	m.Patterns = patterns
	return m
}

func GetMatchCount(list string, applyPatch bool, tailRecursionRule int) int {
	m := NewMatcher(list)
	if applyPatch {
		m.PatchForPartTwo(tailRecursionRule)
	}
	ruleId := 0

	for patternId := range m.Patterns {
		if m.PatternMatchesRule(patternId, ruleId) {
			m.MatchedPatterns = append(m.MatchedPatterns, patternId)
		}
	}
	return len(m.MatchedPatterns)
}

func main() {
	dat, err := ioutil.ReadFile("aoc19.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	num := GetMatchCount(txt, false, 8)
	// Part one
	fmt.Println("Part one:", num)

	num = GetMatchCount(txt, true, 8) // could also put 42...
	// Part one
	fmt.Println("Part two:", num)
}
