package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type AdapterBag struct {
	ResolvedAdapters map[int]int
	Adapters []int // sorted in asc
}

func (bag *AdapterBag) Init(path string) {
	bag.ResolvedAdapters = make(map[int]int)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	adapters := []int{0}
	for _, line := range lines {
		num, err := strconv.Atoi(line)
		if err == nil {
			adapters = append(adapters, num)
		}
	}
	sort.Ints(adapters)
	bag.Adapters = adapters
}

func (bag *AdapterBag) MapJoltDifference() int {
	currentJolt := 0
	distrib := map[int]int{1: 0, 2: 0, 3: 1}
	for _, adapter := range bag.Adapters[1:] {
		diff := adapter - currentJolt
		if diff < 1 || diff > 3 {
			panic(fmt.Sprintf("Invalid joltage diff %d for adapter %d", diff, adapter))
		}
		distrib[diff]++
		currentJolt += diff
	}

	return distrib[1] * distrib[3]
}

func (bag *AdapterBag) GetOptionCount(adapterIndex int) int {
	if paintedNum, exists := bag.ResolvedAdapters[adapterIndex]; exists {
		return paintedNum
	}

	children := bag.GetAdaptersChildren(adapterIndex)

	if len(children) == 0 {
		// Reached the end
		return 1
	}

	numConnections := 0
	for _, child := range children {
		num := bag.GetOptionCount(child)
		bag.ResolveOptionsForAdapter(child, num)
		numConnections += num
	}
	return numConnections
}

func (bag *AdapterBag) ResolveOptionsForAdapter(adapterIndex, options int) {
	bag.ResolvedAdapters[adapterIndex] = options
}

// Get children sorted in descending order
func (bag *AdapterBag) GetAdaptersChildren(adapterIndex int) []int {
	var children []int
	for i := adapterIndex + 1; i < adapterIndex + 4 && i < len(bag.Adapters); i++ {
		if bag.Adapters[i] - bag.Adapters[adapterIndex] > 3 {
			break
		}
		children = append(children, i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(children)))
	return children
}

func main() {
	bag := AdapterBag{}
	bag.Init("aoc10.txt")
	fmt.Println("Part one:", bag.MapJoltDifference())
	fmt.Println("Part two:", bag.GetOptionCount(0))
}
