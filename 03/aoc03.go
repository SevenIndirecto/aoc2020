package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Point struct {
	x int
	y int
}

type TreeMap struct {
	points map[Point]bool
	width  int
	height int
}

func LoadMap(path string) TreeMap {
	points := map[Point]bool{}

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	txt = strings.TrimRight(txt, "\n")
	lines := strings.Split(txt, "\n")

	width := 0
	height := 0
	for y, line := range lines {
		if len(line) < 1 {
			continue
		}
		for x, point := range line {
			points[Point{x, y}] = point == '#'

			if x > width {
				width = x
			}
		}
		height++
	}
	width++
	return TreeMap{points, width, height}
}

func TreesOnSlope(deltaX int, deltaY int, treeMap TreeMap) int {
	loc := Point{0, 0}
	trees := 0

	for loc.y < treeMap.height-1 {
		loc = Point{
			x: (loc.x + deltaX) % treeMap.width,
			y: loc.y + deltaY,
		}

		if treeMap.points[loc] {
			trees++
		}
	}
	return trees
}

func main() {
	treeMap := LoadMap("aoc03.txt")
	trees := TreesOnSlope(3, 1, treeMap)
	fmt.Println("Part one:", trees)

	mul := 1
	slopes := [...][2]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	for _, slope := range slopes {
		mul *= TreesOnSlope(slope[0], slope[1], treeMap)
	}
	fmt.Println("Part two:", mul)
}
