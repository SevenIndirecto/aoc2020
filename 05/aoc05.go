package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func (s Seat) SeatId() int {
	return s.Row*8 + s.Column
}

type Seat struct {
	Row    int
	Column int
}

func ToSeat(pass string) Seat {
	return Seat{
		Row:    getCoordinate(pass[:7], 127, 'F', 'B'),
		Column: getCoordinate(pass[7:10], 7, 'L', 'R'),
	}
}

func getCoordinate(key string, maxSize int, lowKey rune, highKey rune) int {
	min, max := 0, maxSize
	for _, char := range key {
		delta := (max - min + 1) / 2
		if char == highKey {
			min += delta
		} else if char == lowKey {
			max -= delta
		} else {
			panic("Invalid key")
		}
	}
	return min
}

func main() {
	dat, err := ioutil.ReadFile("aoc05.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	var seats []Seat

	for _, line := range lines {
		if len(line) < 10 {
			continue
		}
		seat := ToSeat(line)
		seats = append(seats, seat)
	}

	if len(seats) < 1 {
		panic("Invalids seats")
		return
	}

	sort.Slice(seats, func(i, j int) bool {
		return seats[i].SeatId() > seats[j].SeatId()
	})
	fmt.Println("Part one:", seats[0].SeatId())

	for i := range seats {
		if seats[i].SeatId()-seats[i+1].SeatId() != 1 {
			fmt.Println("Part two:", seats[i].SeatId()-1)
			break
		}
	}
}
