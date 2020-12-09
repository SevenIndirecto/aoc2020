package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type XMAS struct {
	// Keep a count of occurrence of each sum
	PrecomputedSums PrecomputedSums
	Numbers         []int
	PreambleSize    int
	Pos             int
}

type PrecomputedSums map[int]int

func (xmas *XMAS) IncrementPos() {
	// Update count of sums cause by number being removed
	indexPendingRemoval := xmas.Pos - xmas.PreambleSize
	for i := indexPendingRemoval + 1; i < xmas.Pos; i++ {
		sum := xmas.Numbers[i] + xmas.Numbers[indexPendingRemoval]
		if xmas.PrecomputedSums[sum] > 1 {
			xmas.PrecomputedSums[sum]--
		} else {
			delete(xmas.PrecomputedSums, sum)
		}
	}
	xmas.Pos++
	// Precompute sum for new "preamble"
	currentNum := xmas.Numbers[xmas.Pos-1]
	for i := xmas.Pos - xmas.PreambleSize; i < xmas.Pos-1; i++ {
		sum := xmas.Numbers[i] + currentNum
		if _, exists := xmas.PrecomputedSums[sum]; exists {
			xmas.PrecomputedSums[sum]++
		} else {
			xmas.PrecomputedSums[sum] = 1
		}
	}
}

func (xmas *XMAS) Init() {
	xmas.PrecomputedSums = make(PrecomputedSums)
	for i := 0; i < xmas.PreambleSize; i++ {
		for j := i + 1; j < xmas.PreambleSize; j++ {
			sum := xmas.Numbers[i] + xmas.Numbers[j]
			if _, exist := xmas.PrecomputedSums[sum]; exist {
				xmas.PrecomputedSums[sum]++
			} else {
				xmas.PrecomputedSums[sum] = 1
			}
		}
	}
	xmas.Pos = xmas.PreambleSize
}

func (xmas *XMAS) IsCurrentPosValid() bool {
	checksum := xmas.Numbers[xmas.Pos]
	_, exists := xmas.PrecomputedSums[checksum]
	return exists
}

func FindFirstInvalid(path string, preambleSize int) (int, XMAS) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	var numbers []int
	for _, line := range lines {
		num, err := strconv.Atoi(line)
		if err == nil {
			numbers = append(numbers, num)
		}
	}

	xmas := XMAS{
		Numbers:      numbers,
		PreambleSize: preambleSize,
	}
	xmas.Init()

	for valid := true; valid; {
		valid = xmas.IsCurrentPosValid()
		if valid {
			xmas.IncrementPos()
		}
	}
	return xmas.Numbers[xmas.Pos], xmas
}

func FindContiguous(numbers []int, target int) int {
	// Keep a running contiguous slice
	start := 0
	end := 0
	sum := 0

	for ; ; {
		// Increase start until sum < target, and decrease sum as we go
		for i := start; sum > target; i++ {
			sum -= numbers[i]
			start = i + 1
		}

		if sum == target {
			break
		}

		// Increase end until sum > target, and increase sum as we go
		for i := end; sum < target; i++ {
			sum += numbers[i]
			end = i + 1
		}

		if sum == target {
			break
		}
	}

	chunk := numbers[start : end]
	sort.Ints(chunk)
	return chunk[0] + chunk[len(chunk)-1]
}

func main() {
	invalidNum, xmas := FindFirstInvalid("aoc09.txt", 25)
	fmt.Println("Part one:", invalidNum)
	fmt.Println("Part two:", FindContiguous(xmas.Numbers, invalidNum))
}
