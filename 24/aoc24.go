package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Point struct {
	x, y int
}

type TileSet struct {
	tiles map[Point]Color
}

type Color bool

const (
	BLACK = true
	WHITE = false
)

func GetTilePoint(path string, origin Point) Point {
	p := Point{x: origin.x, y: origin.y}

	for i := 0; i < len(path); i++ {
		dir := string(path[i])
		if dir == "s" || dir == "n" {
			i++
			dir += string(path[i])
		}

		switch dir {
		case "e":
			p.x += 2
		case "se":
			p.x += 1
			p.y -= 1
		case "sw":
			p.x -= 1
			p.y -= 1
		case "w":
			p.x -= 2
		case "nw":
			p.x -= 1
			p.y += 1
		case "ne":
			p.x += 1
			p.y += 1
		}
	}
	return p
}

func NewTileSet() TileSet {
	return TileSet{tiles: map[Point]Color{Point{0, 0}: WHITE}}
}

func (ts *TileSet) Paint(lines []string) {
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		pointToPaint := GetTilePoint(line, Point{0, 0})
		if _, exists := ts.tiles[pointToPaint]; exists {
			ts.tiles[pointToPaint] = !ts.tiles[pointToPaint]
		} else {
			ts.tiles[pointToPaint] = BLACK
		}
	}
}

func (ts *TileSet) GetBlackCount() int {
	count := 0
	for _, color := range ts.tiles {
		if color == BLACK {
			count++
		}
	}
	return count
}

func (ts *TileSet) ExecuteDailyPaints(days int) {
	pointsToCheck := [7]Point{}

	for i := 1; i <= days; i++ {
		snapshot := make(map[Point]Color)
		for k, v := range ts.tiles {
			snapshot[k] = v
		}

		for currentPoint, _ := range ts.tiles {

			// Point and neighbors
			deltas := []delta{
				{2, 0},
				{1, -1},
				{-1, -1},
				{-2, 0},
				{-1, 1},
				{1, 1},
			}
			for i, d := range deltas {
				pointsToCheck[i] = Point{currentPoint.x + d.dx, currentPoint.y + d.dy}
			}
			pointsToCheck[6] = currentPoint

			for _, point := range pointsToCheck {
				blacks := ts.GetBlackNeighbors(point)
				var color Color

				if _, exists := ts.tiles[point]; !exists {
					color = WHITE
				} else {
					color = ts.tiles[point]
				}

				if color == BLACK {
					if blacks == 0 || blacks > 2 {
						snapshot[point] = WHITE
					}
				} else if color == WHITE {
					if blacks == 2 {
						snapshot[point] = BLACK
					}
				}
			}
		}

		ts.tiles = snapshot
	}
}

type delta struct {
	dx, dy int
}

func (ts *TileSet) GetBlackNeighbors(p Point) int {
	deltas := []delta{
		{2, 0},
		{1, -1},
		{-1, -1},
		{-2, 0},
		{-1, 1},
		{1, 1},
	}

	black := 0
	for _, d := range deltas {
		n := Point{p.x + d.dx, p.y + d.dy}
		if color, exists := ts.tiles[n]; exists {
			if color == BLACK {
				black++
			}
		}
	}
	return black
}

func main() {
	dat, err := ioutil.ReadFile("aoc24.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	ts := NewTileSet()
	ts.Paint(lines)
	fmt.Println("Part one:", ts.GetBlackCount())

	ts.ExecuteDailyPaints(100)
	fmt.Println("Part two:", ts.GetBlackCount())
}
