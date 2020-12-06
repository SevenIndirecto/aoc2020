package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Group []string

func GetYesToAnyCount(g Group) int {
	yesMap := make(map[rune]struct{})

	for _, personsAnswers := range g {
		for _, answer := range personsAnswers {
			yesMap[answer] = struct{}{}
		}
	}
	return len(yesMap)
}

func GetYesToAllCount(g Group) int {
	yesMap := make(map[rune]int)

	for _, personsAnswers := range g {
		for _, answer := range personsAnswers {
			yesMap[answer] += 1
		}
	}

	allAnsweredYesCount := 0
	for _, answeredYesCount := range yesMap {
		if answeredYesCount == len(g) {
			allAnsweredYesCount++
		}
	}
	return allAnsweredYesCount
}

func main() {
	dat, err := ioutil.ReadFile("aoc06.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	var group Group
	var groups []Group

	for _, line := range lines {
		line = strings.TrimRight(line, "\n")
		if len(line) < 1 {
			groups = append(groups, group)
			group = Group{}
			continue
		}

		group = append(group, line)
	}

	countAny, countAll := 0, 0
	for _, group := range groups {
		countAny += GetYesToAnyCount(group)
		countAll += GetYesToAllCount(group)
	}
	fmt.Println("Part one: ", countAny)
	fmt.Println("Part two: ", countAll)
}
