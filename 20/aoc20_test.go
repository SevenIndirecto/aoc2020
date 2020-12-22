package main

import (
	"fmt"
	"strings"
	"testing"
)


func TestNewTile(t *testing.T) {
	tiles := `
.###
....
...#
#...`
	expectedGrid := `[[0 1 1 1] [0 0 0 0] [0 0 0 1] [1 0 0 0]]`

	tile := NewTile(1, toLines(tiles))
	got := fmt.Sprintf("%v", tile.Grid)

	if got != expectedGrid {
		t.Errorf("Error loading tile got %s expected %s", got, expectedGrid)
	}
}

type FixtureCalcBorderId struct {
	Border BorderType
	Expected int
}

func TestTile_CalculateBorderId(t *testing.T) {
	tileStr := `
.###
....
...#
#...`

	tile := NewTile(1, toLines(tileStr))
	fixtures := []FixtureCalcBorderId{
		{BorderType{TOP, false}, 14},
		{BorderType{RIGHT, false}, 5},
		{BorderType{BOTTOM, false}, 8},
		{BorderType{LEFT, false}, 1},
		{BorderType{TOP, true}, 7},
		{BorderType{RIGHT, true}, 8},
		{BorderType{BOTTOM, true}, 1},
		{BorderType{LEFT, true}, 10},
	}

	for _, f := range fixtures {
		got := tile.CalculateBorderId(f.Border)

		if got != f.Expected {
			t.Errorf("Calculate (%v) got %d expected %d", f.Border, got, f.Expected)
		}
	}
}

func TestTile_PrecomputeBorders(t *testing.T) {
	tileStr1 := `
.###
....
...#
#...`
	tileStr2 := `
###.
...#
...#
#...`

	tile1 := NewTile(1, toLines(tileStr1))
	tile2 := NewTile(2, toLines(tileStr2))
	borderMap := NewPrecomputedBorders()
	tile1.PrecomputeBorders(&borderMap)
	tile2.PrecomputeBorders(&borderMap)
	got := len(borderMap.Map)
	expected := 8
	if got != expected {
		t.Errorf("Failed to precompute borders got %d expected %d", got, expected)
	}

	gotTile1 := fmt.Sprintf("%v", tile1.BorderIds)
	expectedIds := `map[{0 false}:14 {0 true}:7 {1 false}:5 {1 true}:8 {2 false}:8 {2 true}:1 {3 false}:1 {3 true}:10]`
	if gotTile1 != expectedIds {
		t.Errorf("Failed to precompute tile borders got %s expected %s", gotTile1, expectedIds)
	}
}

func toLines(tilesStr string) []string {
	lines := strings.Split(tilesStr, "\n")
	return lines[1:]
}

func TestNewSolver(t *testing.T) {
	solver := NewSolver("aoc20_test1.txt")
	ct := solver.GetCornerTiles()
	mul := 1
	for k, _ := range ct {
		mul *= k
	}
	expected := 20899048083289

	if mul != expected {
		t.Errorf("Failed to match border corners got %d, expected %d", mul, expected)
	}
}

func TestTile_RotateClockwise(t *testing.T) {
	tileStr := `
.###
....
...#
#...`
	expected := `#...
...#
...#
.#.#
`

	tile := NewTile(1, toLines(tileStr))
	tile.RotateClockwise()

	got := tile.String()

	if got != expected {
		t.Errorf("Failed to rotate got \n\n%s\n\nexpected\n\n%s", got, expected)
	}
}

func TestTile_FlipHorizontal(t *testing.T) {
	tileStr := `
.###
....
...#
#...`
	expected := `###.
....
#...
...#
`

	tile := NewTile(1, toLines(tileStr))
	tile.FlipHorizontal()

	got := tile.String()

	if got != expected {
		t.Errorf("Failed to flip got \n\n%s\n\nexpected\n\n%s", got, expected)
	}
}

func TestSolver_ConstructImage(t *testing.T) {
	solver := NewSolver("aoc20_test1.txt")
	img := solver.ConstructImage()
	size := len(img)
	got := img[0][0].Id * img[0][size-1].Id * img[size-1][size-1].Id * img[size-1][0].Id
	expected := 20899048083289

	if got != expected {
		t.Errorf("Failed to consturct image got %d expected %d", got, expected)
	}
}


func TestTile_CalibrateAndMarkMonsters(t *testing.T) {
	solver := NewSolver("aoc20_test1.txt")
	tilePatch := solver.ConstructImage()
	img := TilePatchToImage(tilePatch)
	img.CalibrateAndMarkMonsters()
	expected := 273
	got := img.GetRoughCount()

	if got != expected {
		t.Errorf("Could not mark monsters got %d expected %d rough tiles", got, expected)
	}
}
