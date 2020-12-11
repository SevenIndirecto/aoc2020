package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type PointState string

const (
	Occupied PointState = "#"
	Floor               = "."
	Empty               = "L"
)

type Mode int

const (
	PartOneMode Mode = 0
	PartTwoMode Mode = 1
)

type Point struct {
	X int
	Y int
}

type SeatLayout struct {
	Grid     [][]PointState
	Snapshot [][]PointState
	Mode     Mode
}

func (sl *SeatLayout) Init(path string) {
	sl.Grid = [][]PointState{}
	sl.Snapshot = [][]PointState{}

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		var row []PointState
		for _, point := range line {
			row = append(row, PointState(point))
		}
		sl.Grid = append(sl.Grid, row)
	}
}

func (sl *SeatLayout) PrintState() {
	for _, gridRow := range sl.Grid {
		row := ""
		for _, state := range gridRow {
			row += string(state)
		}
		fmt.Printf("[%s]\n", row)
	}
	fmt.Println()
}

func (sl *SeatLayout) Tick() {
	sl.Snapshot = make([][]PointState, len(sl.Grid))
	for i := range sl.Grid {
		sl.Snapshot[i] = make([]PointState, len(sl.Grid[i]))
		copy(sl.Snapshot[i], sl.Grid[i])
	}

	var adjacentState map[PointState]int

	for y, row := range sl.Snapshot {
		for x := range row {
			point := Point{x, y}

			occupiedThreshold := 4
			if sl.Mode == PartOneMode {
				occupiedThreshold = 4
				adjacentState = GetAdjacentState(sl.Snapshot, point)
			} else {
				occupiedThreshold = 5
				adjacentState = GetFirstVisibleChairs(sl.Snapshot, point)
			}
			sl.Grid[y][x] = GetNewValue(adjacentState, sl.Snapshot, point, occupiedThreshold)
		}
	}
}

func (sl *SeatLayout) StateMap() map[PointState]int {
	m := map[PointState]int{Occupied: 0, Floor: 0, Empty: 0}
	for _, row := range sl.Grid {
		for _, point := range row {
			m[point]++
		}
	}
	return m
}

func GetNewValue(compareState map[PointState]int, grid [][]PointState, point Point, threshold int) PointState {
	if grid[point.Y][point.X] == Empty {
		if compareState[Occupied] == 0 {
			return Occupied
		}
	} else if grid[point.Y][point.X] == Occupied {
		if compareState[Occupied] >= threshold {
			return Empty
		}
	}
	return grid[point.Y][point.X]
}

func GetAdjacent(grid [][]PointState, point Point) []Point {
	var adjacent []Point

	for ny := point.Y - 1; ny <= point.Y+1; ny++ {
		if ny >= len(grid) || ny < 0 {
			continue
		}

		for nx := point.X - 1; nx <= point.X+1; nx++ {
			if nx >= len(grid[ny]) || nx < 0 || (nx == point.X && ny == point.Y) {
				continue
			}

			adjacent = append(adjacent, Point{X: nx, Y: ny})
		}
	}
	return adjacent
}

func GetAdjacentState(grid [][]PointState, point Point) map[PointState]int {
	pointMap := map[PointState]int{Empty: 0, Floor: 0, Occupied: 0}
	adjacent := GetAdjacent(grid, point)
	for _, a := range adjacent {
		pointMap[grid[a.Y][a.X]]++
	}
	return pointMap
}

func GetFirstVisibleChairs(grid [][]PointState, point Point) map[PointState]int {
	pointMap := map[PointState]int{Empty: 0, Floor: 0, Occupied: 0}
	adjacent := GetAdjacent(grid, point)

	for _, a := range adjacent {
		// Set direction
		dirX := a.X - point.X
		dirY := a.Y - point.Y

		// Get first visible chair
		candidate := Point{X: point.X, Y: point.Y}
		for ; ; {
			candidate = Point{X: candidate.X + dirX, Y: candidate.Y + dirY}
			if candidate.Y < 0 || candidate.X < 0 || candidate.Y >= len(grid) || candidate.X >= len(grid[0]) {
				break
			}
			t := grid[candidate.Y][candidate.X]
			if t == Occupied || t == Empty {
				pointMap[t]++
				break
			}
		}
	}

	return pointMap
}

func FindEquilibrium(sl *SeatLayout) int {
	for ; ; {
		sl.Tick()
		allEqual := true
		for i, _ := range sl.Grid {
			if !reflect.DeepEqual(sl.Grid[i], sl.Snapshot[i]) {
				allEqual = false
				break
			}
		}
		if allEqual {
			break
		}
	}
	return sl.StateMap()[Occupied]
}

func main() {
	sl := SeatLayout{Mode: PartOneMode}
	sl.Init("aoc11.txt")
	fmt.Println("Part one: ", FindEquilibrium(&sl))

	sl = SeatLayout{Mode: PartTwoMode}
	sl.Init("aoc11.txt")
	fmt.Println("Part two: ", FindEquilibrium(&sl))
}
