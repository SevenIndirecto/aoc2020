package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type NumberTrail struct {
	Latest int
	Prev int
}

func (nt *NumberTrail) Age() int {
	if nt.Latest == 0 || nt.Prev == 0 {
		return 0
	}
	return nt.Latest - nt.Prev
}

func (nt *NumberTrail) Update(turn int) {
	if nt.Latest == 0 {
		nt.Latest = turn
	} else {
		nt.Prev = nt.Latest
		nt.Latest = turn
	}
}

type Memorizer struct {
	Memory map[int]*NumberTrail
	Turn int
	LastNumberSpoken int
}

func NewMemorizer(numbers []int) Memorizer {
	m := Memorizer{Memory: make(map[int]*NumberTrail)}
	for _, number := range numbers {
		m.Turn++
		m.AddNewNumber(number, m.Turn)
		m.LastNumberSpoken = number
	}
	return m
}

func (m *Memorizer) AddNewNumber(num, turn int) {
	m.Memory[num] = &NumberTrail{Latest: turn, Prev: 0}
}

func (m *Memorizer) GetNthNumberSpoken(target int) int {
	for m.Turn < target {
		m.Turn++

		newNumberSpoken := m.Memory[m.LastNumberSpoken].Age()

		if _, found := m.Memory[newNumberSpoken]; !found {
			m.AddNewNumber(newNumberSpoken, m.Turn)
		} else {
			m.Memory[newNumberSpoken].Update(m.Turn)
		}
		m.LastNumberSpoken = newNumberSpoken
	}
	return m.LastNumberSpoken
}

func main() {
	dat, err := ioutil.ReadFile("aoc15.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	numStrings := strings.Split(lines[0], ",")
	var nums []int
	for _, nStr := range numStrings {
		n, _ := strconv.Atoi(nStr)
		nums = append(nums, n)
	}

	m := NewMemorizer(nums)
	fmt.Println("Part one:", m.GetNthNumberSpoken(2020))
	fmt.Println("Part two:", m.GetNthNumberSpoken(30000000))
}
