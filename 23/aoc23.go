package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type CupGame struct {
	// This acts as a map, the index represents a Cup (label) and it points to the next Cup (the value it contains)
	// -> Had to get this idea of the net, but really like it, more efficient than using map[int]int
	cups []int
	current int
	max int
}

func NewCupGame(seed string, max int) CupGame {
	cg := CupGame{cups: make([]int, max+1), current: 0, max: max}
	start := 0

	// Initialize
	for _, digit := range seed {
		label, _ := strconv.Atoi(string(digit))
		if start == 0 {
			start = label
		}

		cg.cups[cg.current] = label
		cg.current = label
	}
	// For part two
	for i := len(seed)+1; i <= max; i++ {
		cg.cups[cg.current] = i
		cg.current = i
	}
	// Connect circle
	cg.cups[cg.current] = cg.cups[0]

	cg.current = start

	return cg
}

func (cg *CupGame) Play(rounds int) {
	for {
		cup1 := cg.cups[cg.current]
		cup2 := cg.cups[cup1]
		cup3 := cg.cups[cup2]
		after := cg.cups[cup3]

		cg.cups[cg.current] = after

		// Find destination
		dest := cg.current - 1

		for {
			if dest == 0 {
				dest = cg.max
			}

			if dest == cup1 || dest == cup2 || dest == cup3 {
				dest--
			} else {
				break
			}
		}

		cg.cups[cup3] = cg.cups[dest]
		cg.cups[dest] = cup1

		// Next
		cg.current = cg.cups[cg.current]

		rounds--
		if rounds == 0 {
			break
		}
	}
}

func (cg *CupGame) Print(max int) {
	str := ""
	current := cg.current

	for i := 0; i < cg.max; i++ {
		str += fmt.Sprintf("%d ", current)
		current = cg.cups[current]

		if i > max {
			break
		}
	}

	fmt.Println("Cups:", str)
}

func (cg *CupGame) GetPartOneSig() string {
	str := ""
	current := cg.cups[1]
	for ; current != 1; current = cg.cups[current] {
		str += strconv.Itoa(current)
	}
	return str
}

func (cg *CupGame) GetPartTwoSig() int {
	return cg.cups[1] * cg.cups[cg.cups[1]]
}

func main() {
	dat, err := ioutil.ReadFile("aoc23.txt")
	if err != nil {
		panic(err)
	}

	seed := strings.TrimRight(string(dat), "\n")

	cg := NewCupGame(seed, 9)
	cg.Play(100)
	fmt.Println("Part one:", cg.GetPartOneSig())

	cg = NewCupGame(seed, 1000000)
	cg.Play(10000000)
	fmt.Println("Part two:", cg.GetPartTwoSig())
}
