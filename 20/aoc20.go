package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	TOP = iota
	RIGHT
	BOTTOM
	LEFT
)

type Solver struct {
	Precomputed PrecomputedBorders
	Tiles       map[int]Tile
}

type PrecomputedBorders struct {
	Map map[int][]Border
}

func NewPrecomputedBorders() PrecomputedBorders {
	return PrecomputedBorders{
		Map: make(map[int][]Border),
	}
}

type Border struct {
	Id     int
	TileId int
	Type   BorderType
}

type BorderType struct {
	Side    int
	Flipped bool
}

type Tile struct {
	Id          int
	Grid        [][]int
	FlippedGrid [][]int
	BorderIds   map[BorderType]int
	Rotation    int
	Flipped     bool
}

func (tile Tile) ToStringRows() []string {
	var stringRows []string

	for y := range tile.Grid {
		rowStr := ""
		for x := range tile.Grid[y] {
			painted := tile.Grid[y][x]
			if painted == 2 {
				rowStr += "O"
			} else if painted == 1 {
				rowStr += "#"
			} else {
				rowStr += "."
			}
		}
		stringRows = append(stringRows, rowStr)
	}
	return stringRows
}

func (tile Tile) String() string {
	str := ""
	for _, row := range tile.ToStringRows() {
		str += fmt.Sprintf("%s\n", row)
	}
	return str
}

func (tile *Tile) RotateClockwise() {
	size := len(tile.Grid)
	rotatedGrid := make([][]int, size)
	for y := range rotatedGrid {
		rotatedGrid[y] = make([]int, size)
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			ry := x
			rx := size - 1 - y
			rotatedGrid[ry][rx] = tile.Grid[y][x]
		}
	}
	tile.Grid = rotatedGrid
}

func (tile *Tile) FlipHorizontal() {
	size := len(tile.Grid)
	flippedGrid := make([][]int, size)
	for y := range flippedGrid {
		flippedGrid[y] = make([]int, size)
	}


	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			flippedGrid[y][size - 1 - x] = tile.Grid[y][x]
		}
	}
	tile.Grid = flippedGrid
}

func (tile Tile) PrecomputeBorders(borderMap *PrecomputedBorders) {
	variants := []BorderType{
		{TOP, false},
		{LEFT, false},
		{BOTTOM, false},
		{RIGHT, false},
		{TOP, true},
		{LEFT, true},
		{BOTTOM, true},
		{RIGHT, true},
	}

	for _, borderType := range variants {
		border := Border{
			Id:     tile.CalculateBorderId(borderType),
			TileId: tile.Id,
			Type:   borderType,
		}
		tile.BorderIds[borderType] = border.Id

		if _, exists := borderMap.Map[border.Id]; !exists {
			borderMap.Map[border.Id] = []Border{border}
		} else {
			borderMap.Map[border.Id] = append(borderMap.Map[border.Id], border)
		}
	}
}

func p2int(power int) int {
	return int(math.Pow(2, float64(power)))
}

func ModPositive(a, divisor int) int {
	return ((a % divisor) + divisor) % divisor
}

// Flipped assumes horizontal flip
func (tile Tile) CalculateBorderId(borderType BorderType) int {
	id := 0
	switch borderType.Side {
	case TOP:
		row := tile.Grid[0]
		for x := 0; x < len(row); x++ {
			if borderType.Flipped {
				id += row[len(row)-1-x] * p2int(x)
			} else {
				id += row[x] * p2int(x)
			}
		}
	case RIGHT:
		x := len(tile.Grid[0]) - 1
		if borderType.Flipped {
			x = 0
		}
		for y := 0; y < len(tile.Grid); y++ {
			id += tile.Grid[y][x] * p2int(y)
		}
	case BOTTOM:
		row := tile.Grid[len(tile.Grid)-1]
		for x := 0; x < len(row); x++ {
			if borderType.Flipped {
				id += row[x] * p2int(x)
			} else {
				id += row[len(row)-1-x] * p2int(x)
			}
		}
	case LEFT:
		x := 0
		if borderType.Flipped {
			x = len(tile.Grid[0]) - 1
		}
		for y := 0; y < len(tile.Grid); y++ {
			id += tile.Grid[len(tile.Grid)-1-y][x] * p2int(y)
		}
	}
	return id
}

func NewTile(id int, gridLines []string) Tile {
	tile := Tile{
		Id:        id,
		BorderIds: make(map[BorderType]int),
	}
	var grid [][]int

	for _, line := range gridLines {
		if len(line) < 2 {
			continue
		}
		var row []int

		for _, c := range line {
			painted := 0
			if string(c) == "#" {
				painted = 1
			}
			row = append(row, painted)
		}
		grid = append(grid, row)
	}
	tile.Grid = grid
	return tile
}

func NewSolver(path string) Solver {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	solver := Solver{
		Precomputed: NewPrecomputedBorders(),
		Tiles:       make(map[int]Tile),
	}

	var tileLines []string
	tileId := -1
	re := regexp.MustCompile("Tile ([0-9]+):")

	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		if match != nil {
			// Start new
			tileId, _ = strconv.Atoi(match[1])
			continue
		}

		if len(line) < 2 && len(tileLines) > 0 {
			// End of tile definition
			tile := NewTile(tileId, tileLines)
			tile.PrecomputeBorders(&solver.Precomputed)
			solver.Tiles[tile.Id] = tile
			tileLines = make([]string, 0)
			continue
		}
		tileLines = append(tileLines, line)
	}
	return solver
}

// Get map of Tiles which are considered as border tiles, due to having a unique borderId
func (solver Solver) GetBorderTiles() map[int][]Border {
	tileIdsToBorders := make(map[int][]Border)

	for _, borders := range solver.Precomputed.Map {
		if len(borders) < 2 {
			tileId := borders[0].TileId
			if _, exists := tileIdsToBorders[tileId]; !exists {
				tileIdsToBorders[tileId] = []Border{}
			}
			tileIdsToBorders[tileId] = append(tileIdsToBorders[tileId], borders[0])
		}
	}
	return tileIdsToBorders
}

// A corner tile has exactly 4 unique border IDs, two for each corner border and 2 flipped
func (solver Solver) GetCornerTiles() map[int][]Border {
	borderTiles := solver.GetBorderTiles()
	cornerTiles := make(map[int][]Border)
	for k, v := range borderTiles {
		if len(v) == 4 {
			cornerTiles[k] = v
		}
	}
	return cornerTiles
}

func (solver Solver) ConstructImage() [][]Tile {
	squareSize := int(math.Sqrt(float64(len(solver.Tiles))))
	fmt.Printf("Building %dx%d image\n", squareSize, squareSize)

	ct := solver.GetCornerTiles()
	//Find top left Corner
	topLeftTileId := -1
	for tileId, borders := range ct {
		matchCount := 0
		for _, borderType := range borders {
			if borderType.Type.Side == LEFT && !borderType.Type.Flipped ||
				borderType.Type.Side == TOP && !borderType.Type.Flipped {
				matchCount++
			}
		}

		if matchCount == 2 {
			topLeftTileId = tileId
			break
		}
	}

	// Init
	s := make([][]Tile, squareSize)
	for y := range s {
		s[y] = make([]Tile, squareSize)
	}
	sidesToCheck := [2][3]int{
		{0, 1, RIGHT},
		{1, 0, BOTTOM},
	}
	usedTiles := make(map[int]bool)

	for y := 0; y < squareSize; y++ {
		for x := 0; x < squareSize; x++ {
			if x == 0 && y == 0 {
				// bootstrap by seeding with first tile
				s[0][0] = solver.Tiles[topLeftTileId]
				usedTiles[s[0][0].Id] = true
			}

			currentTile := s[y][x]

			// Find matching tile for each side (Only need bottom and right really)
			for _, side := range sidesToCheck {
				dy, dx, targetSide := side[0], side[1], side[2]
				ny, nx := y + dy, x + dx
				if ny < 0 || nx < 0 || ny >= squareSize || nx >= squareSize {
					// Out of bounds
					continue
				}

				if s[ny][nx].Id > 0 {
					// Already inserted -> skip
					continue
				}

				// Could use precomputed but meh for now...
				var targetBorderId int
				var targetMachingSide int
				// Only looking right and bottom...
				if targetSide == RIGHT {
					targetMachingSide = LEFT
					targetBorderId = currentTile.CalculateBorderId(BorderType{LEFT, true})
				} else if targetSide == BOTTOM {
					targetMachingSide = TOP
					targetBorderId = currentTile.CalculateBorderId(BorderType{BOTTOM, true})
				}

				// Find matching tile
				var matchingBorder Border
				candidates := solver.Precomputed.Map[targetBorderId]
				if len(candidates) != 2 {
					panic("Did not get 2 candidates as expected")
				}
				for _, c := range candidates {
					if c.TileId != currentTile.Id {
						matchingBorder = c
						break
					}
				}

				// Calibrate matching tile and insert
				matchingTile := solver.Tiles[matchingBorder.TileId]
				// 1. Flip if needed
				if matchingBorder.Type.Flipped {
					matchingTile.FlipHorizontal()
				}
				// 2. handle rotation
				numRotations := ModPositive(targetMachingSide - matchingBorder.Type.Side, 4)
				for i := 0; i < numRotations; i++ {
					matchingTile.RotateClockwise()
				}

				if _, exists := usedTiles[matchingTile.Id]; exists {
					panic(fmt.Sprintf("Tile %v already used", matchingTile))
				}
				s[ny][nx] = matchingTile
				usedTiles[matchingTile.Id] = true
			}
		}
	}

	PrintTilePatch(s)
	return s
}

func PrintTilePatch(tiles [][]Tile) {
	for y := 0; y < len(tiles); y++ {
		if tiles[y][0].Id == 0 {
			// Skip row of unset tiles
			continue
		}
		linesInRow := make([][]string, 0)

		for x := 0; x < len(tiles[y]); x++ {
			tile := tiles[y][x]
			lines := tile.ToStringRows()
			if len(lines) > 0 {
				linesInRow = append(linesInRow, lines)
			}
		}

		for lineY := 0; lineY < len(linesInRow[0]); lineY++ {
			row := ""
			for i := 0; i < len(linesInRow); i++ {
				row += linesInRow[i][lineY] + " "
			}
			fmt.Println(row)
		}
		fmt.Println()
	}
}

func TilePatchToImage(tiles [][]Tile) Tile {
	// Remove border from each tile
	borderlessTileGridSize := len(tiles[0][0].Grid) - 2
	grid := make([][]int, len(tiles) * borderlessTileGridSize)
	for y := range grid {
		grid[y] = make([]int, len(tiles) * borderlessTileGridSize)
	}

	for tileY := range tiles {
		for tileX, tile := range tiles[tileY] {

			for gridY := 1; gridY < len(tile.Grid) - 1; gridY++ {
				for gridX := 1; gridX < len(tile.Grid[gridY]) - 1; gridX++ {
					y := tileY * borderlessTileGridSize + gridY - 1
					x := tileX * borderlessTileGridSize + gridX - 1
					grid[y][x] = tile.Grid[gridY][gridX]
				}
			}
		}
	}
	return Tile{Grid: grid}
}

type Delta struct {
	dy, dx int
}

func (tile *Tile) CalibrateAndMarkMonsters() {
	g := tile.Grid

	deltasToCheck := []Delta{
		{0, 0},
		{1, 1},
		{1, 4},
		{0, 5},
		{0, 6},
		{1, 7},
		{1, 10},
		{0, 11},
		{0, 12},
		{1, 13},
		{1, 16},
		{0, 17},
		{0, 18},
		{0, 19},
		{-1, 18},
	}

	usedFlip := false
	rotateCount := 0
	for ; ; {
		roughCountBefore := tile.GetRoughCount()

		for y := range g {
			for x := range g[y] {
				found := true
				// Find pattern
				for _, delta := range deltasToCheck {
					nx := x + delta.dx
					ny := y + delta.dy

					if nx < 0 || ny < 0 || nx >= len(g) || ny >= len(g) || g[ny][nx] != 1 {
						found = false
						break
					}
				}

				if found {
					// Mark tiles with monster
					for _, delta := range deltasToCheck {
						nx := x + delta.dx
						ny := y + delta.dy
						g[ny][nx] = 2
					}
				}
			}
		}

		roughCountAfter := tile.GetRoughCount()
		if roughCountAfter != roughCountBefore {
			return // Done
		}

		// Try to reorient
		if rotateCount >= 3 {
			if usedFlip {
				panic("No monsters could be found in any orientation")
			}
			fmt.Println("Recalibrating by flipping horizontally")
			tile.FlipHorizontal()
			usedFlip = true
			rotateCount = 0
		} else {
			fmt.Println("Recalibrating by rotating CW")
			tile.RotateClockwise()
			rotateCount++
		}
		g = tile.Grid
	}
}

func (tile Tile) GetRoughCount() int {
	sum := 0
	for y := range tile.Grid {
		for _, val := range tile.Grid[y] {
			if val == 1 {
				sum++
			}
		}
	}
	return sum
}

func main() {
	solver := NewSolver("aoc20.txt")
	tilePatch := solver.ConstructImage()
	size := len(tilePatch)
	partOne := tilePatch[0][0].Id * tilePatch[0][size-1].Id * tilePatch[size-1][size-1].Id * tilePatch[size-1][0].Id

	fmt.Println("Part one:", partOne)

	img := TilePatchToImage(tilePatch)
	img.CalibrateAndMarkMonsters()

	fmt.Println("Marked monsters")
	fmt.Println(img)

	fmt.Println("Part two:", img.GetRoughCount())
}
