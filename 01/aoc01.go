package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func ExpenseReport(entries []int, target int) (int, int) {
	for i, val := range entries {
		for j := len(entries) - 1; j > i; j-- {
			sum := entries[j] + val
			if sum == target {
				return entries[j], val
			}
			if sum > target {
				break
			}
		}
	}
	return -1, -1
}

func ExpenseReportThree(entries []int, target int) (int, int, int) {
	for _, a := range entries {
		subTarget := target - a
		b, c := ExpenseReport(entries, subTarget)

		if a + b + c == target {
			return a, b, c
		}
	}
	panic("Could not find 2020")
}

func main() {
	dat, err := ioutil.ReadFile("aoc01.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	var entries []int
	for _, i := range lines {
		num, _ := strconv.Atoi(i)
		if num == 0 {
			continue
		}
		entries = append(entries, num)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(entries)))

	a, b := ExpenseReport(entries, 2020)
	fmt.Println("Part 1:", a, b, a * b)
	a, b, c  := ExpenseReportThree(entries, 2020)
	fmt.Println("Part 2:", a, b, c, a * b * c)
}
