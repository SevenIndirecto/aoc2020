package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Policy struct {
	RuleA int
	RuleB int
	Char  rune
}

type DbEntry struct {
	Policy Policy
	Pass   string
}

func IsValidPartOne(entry DbEntry) bool {
	total := 0
	for _, char := range entry.Pass {
		if char == entry.Policy.Char {
			total++
		}
	}
	return total >= entry.Policy.RuleA && total <= entry.Policy.RuleB
}

func IsValidPartTwo(entry DbEntry) bool {
	pass := []rune(entry.Pass)
	indexA := entry.Policy.RuleA - 1
	indexB := entry.Policy.RuleB - 1
	if indexA >= len(pass) || indexB >= len(pass) {
		return false
	}
	return (pass[indexA] == entry.Policy.Char) != (pass[indexB] == entry.Policy.Char)
}

func ParseLine(line string) (DbEntry, error) {
	s := strings.Split(line, ": ")
	if len(s) < 2 {
		return DbEntry{}, errors.New("empty")
	}
	pass := s[1]
	policySplit := strings.Split(s[0], " ")
	char := []rune(policySplit[1])[0]
	ruleSplit := strings.Split(policySplit[0], "-")
	min, _ := strconv.Atoi(ruleSplit[0])
	max, _ := strconv.Atoi(ruleSplit[1])

	policy := Policy{
		RuleA:  min,
		RuleB:  max,
		Char: char,
	}

	return DbEntry{
		Policy: policy,
		Pass:   pass,
	}, nil
}

func main() {
	dat, err := ioutil.ReadFile("aoc02.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	var entries []DbEntry

	for _, line := range lines {
		entry, err := ParseLine(line)
		if nil == err {
			entries = append(entries, entry)
		}
	}
	partOneCount := 0
	partTwoCount := 0
	for _, entry := range entries {
		if IsValidPartOne(entry) {
			partOneCount++
		}
		if IsValidPartTwo(entry) {
			partTwoCount++
		}
	}
	fmt.Println("Part One: ", partOneCount)
	fmt.Println("Part Two: ", partTwoCount)
}
