package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Bus struct {
	Id int
	Offset int
}

func GetFirstBus(buses []int, start int) (int, int, error) {
	earliestBus := -1
	lowestWait := math.MaxInt64

	for _, bus := range buses {
		firstArrivalAfterStart := bus + bus * int(math.Ceil(float64(start / bus)))
		wait := firstArrivalAfterStart - start
		if wait < lowestWait {
			earliestBus = bus
			lowestWait = wait
		}
	}
	if earliestBus == -1 {
		return earliestBus, lowestWait, errors.New("could not find latest bus")
	}
	return earliestBus, lowestWait, nil
}

func FindNextSyncTime(period, start, a, b, offset int) int {
	for t := start; ; t += period {
		if period > 1 {
			if (t + offset) % b == 0 {
				return t
			}
		} else {
			if t % a == 0 && (t + offset) % b == 0 {
				return t
			}
		}
	}
}

func FindFirstAndPeriod(period, a, b, offset int) (int, int) {
	first := FindNextSyncTime(period, a, a, b, offset)
	second := FindNextSyncTime(period, first + period, a, b, offset)
	newPeriod := second - first
	return first, newPeriod
}

func SolveSchedule(buses []Bus) int {
	period := 1
	first := buses[0].Id

	for _, bus := range buses[1:] {
		first, period = FindFirstAndPeriod(period, first, bus.Id, bus.Offset)
	}
	return first
}

func main() {
	dat, err := ioutil.ReadFile("aoc13.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	start, err := strconv.Atoi(lines[0])
	if err != nil {
		panic("Fail")
	}
	splitSchedule := strings.Split(lines[1], ",")
	var buses []int
	for _, bus := range splitSchedule {
		if bus != "x" {
			id, err := strconv.Atoi(bus)
			if err != nil {
				panic("Parse id fail")
			}
			buses = append(buses, id)
		}
	}

	earliestBus, lowestWait, err := GetFirstBus(buses, start)
	if err != nil {
		panic("Could not find any matching bus")
	}
	mul := lowestWait * earliestBus
	fmt.Println("Part one:", mul)

	// Part 2
	var partTwoBuses []Bus

	for i, bus := range splitSchedule {
		if bus == "x" {
			continue
		}

		id, err := strconv.Atoi(bus)
		if err != nil {
			panic("Parse id fail")
		}
		partTwoBuses = append(partTwoBuses, Bus{id, i})
	}
	fmt.Println("Part two: ", SolveSchedule(partTwoBuses))
}
