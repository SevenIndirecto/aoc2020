package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	ACTIVE_REPR   = "#"
	INACTIVE_REPR = "."
)

const (
	ACTIVE   = true
	INACTIVE = false
)

const (
	MODE_PART_ONE = 0
	MODE_PART_TWO = 1
)

type Point struct {
	X, Y, Z, W int
}

type Constraint struct {
	Min, Max int
}

type PocketDimension struct {
	Cycle       int
	Cubes       map[Point]bool
	Snapshot    map[Point]bool
	ConstraintX Constraint
	ConstraintY Constraint
	ConstraintZ Constraint
	ConstraintW Constraint
	Mode        int
}

func NewPocketDimension(pattern string) PocketDimension {
	lines := strings.Split(pattern, "\n")
	cubes := make(map[Point]bool)
	pd := PocketDimension{
		Cycle:       0,
		ConstraintX: Constraint{0, 0},
		ConstraintY: Constraint{0, 0},
		ConstraintZ: Constraint{0, 0},
		ConstraintW: Constraint{0, 0},
		Mode:        MODE_PART_ONE,
	}

	for y, line := range lines {
		if len(line) < 2 {
			continue
		}
		for x, char := range line {
			point := Point{X: x, Y: y, Z: 0, W: 0}
			cubes[point] = string(char) == ACTIVE_REPR
			pd.UpdateConstraints(point)
		}
	}
	pd.Cubes = cubes
	return pd
}

func (pd *PocketDimension) UpdateConstraints(p Point) {
	if p.X < pd.ConstraintX.Min {
		pd.ConstraintX.Min = p.X
	}
	if p.X > pd.ConstraintX.Max {
		pd.ConstraintX.Max = p.X
	}
	if p.Y < pd.ConstraintY.Min {
		pd.ConstraintY.Min = p.Y
	}
	if p.Y > pd.ConstraintY.Max {
		pd.ConstraintY.Max = p.Y
	}
	if p.Z < pd.ConstraintZ.Min {
		pd.ConstraintZ.Min = p.Z
	}
	if p.Z > pd.ConstraintZ.Max {
		pd.ConstraintZ.Max = p.Z
	}
	if p.W < pd.ConstraintW.Min {
		pd.ConstraintW.Min = p.W
	}
	if p.W > pd.ConstraintW.Max {
		pd.ConstraintW.Max = p.W
	}
}

func (pd *PocketDimension) GetNeighborsAndSelf(p Point) [81]Point {
	var neighbors [81]Point

	i := 0

	for deltaX := -1; deltaX <= 1; deltaX++ {
		for deltaY := -1; deltaY <= 1; deltaY++ {
			for deltaZ := -1; deltaZ <= 1; deltaZ++ {
				for deltaW := -1; deltaW <= 1; deltaW++ {
					nx, ny, nz, nw := p.X-deltaX, p.Y-deltaY, p.Z-deltaZ, p.W-deltaW
					neighbors[i] = Point{nx, ny, nz, nw}
					i++
				}
			}
		}
	}
	return neighbors
}

func (pd *PocketDimension) GetNeighborState(p Point) map[bool]int {
	state := map[bool]int{ACTIVE: 0, INACTIVE: 0}

	for deltaX := -1; deltaX <= 1; deltaX++ {
		for deltaY := -1; deltaY <= 1; deltaY++ {
			for deltaZ := -1; deltaZ <= 1; deltaZ++ {
				for deltaW := -1; deltaW <= 1; deltaW++ {
					if deltaX == 0 && deltaY == 0 && deltaZ == 0 && deltaW == 0 {
						continue // skip self
					}
					nx, ny, nz, nw := p.X-deltaX, p.Y-deltaY, p.Z-deltaZ, p.W-deltaW
					neighbor := Point{nx, ny, nz, nw}

					isActive, exists := pd.Snapshot[neighbor]
					if !exists || !isActive {
						state[INACTIVE]++
					} else {
						state[ACTIVE]++
					}
				}
			}
		}
	}
	return state
}

func (pd *PocketDimension) ExecuteCycle() {
	pd.Cycle++
	pd.CreateSnapshot()

	for point := range pd.Snapshot {
		neighborsAndSelf := pd.GetNeighborsAndSelf(point)

		for _, pointToUpdate := range neighborsAndSelf {
			neighborState := pd.GetNeighborState(pointToUpdate)
			wasPreviouslyActive, pointToUpdateExists := pd.Snapshot[pointToUpdate]
			if !pointToUpdateExists {
				wasPreviouslyActive = false
			}

			isNowActive := GetNewState(wasPreviouslyActive, neighborState)

			if !pointToUpdateExists && isNowActive || pointToUpdateExists && isNowActive != wasPreviouslyActive {
				if pd.Mode == MODE_PART_ONE && pointToUpdate.W != 0 {
					continue
				}

				pd.Cubes[pointToUpdate] = isNowActive

				if !pointToUpdateExists {
					pd.UpdateConstraints(pointToUpdate)
				}
			}
		}
	}
}

func (pd *PocketDimension) CreateSnapshot() {
	pd.Snapshot = make(map[Point]bool)

	for point, isActive := range pd.Cubes {
		pd.Snapshot[point] = isActive
	}
}

func (pd *PocketDimension) PrintOut() {
	fmt.Println("MAP", pd.ConstraintX, pd.ConstraintY, pd.ConstraintZ)
	if pd.Cycle == 0 {
		fmt.Println("\nBefore any cycles:")
	} else {
		fmt.Printf("After %d cycle:\n", pd.Cycle)
	}
	for w := pd.ConstraintW.Min; w <= pd.ConstraintW.Max; w++ {
		for z := pd.ConstraintZ.Min; z <= pd.ConstraintZ.Max; z++ {
			fmt.Printf("\nz=%d, w=%d\n\n", z, w)

			for y := pd.ConstraintY.Min; y <= pd.ConstraintY.Max; y++ {
				row := fmt.Sprintf("%d", y)
				for x := pd.ConstraintX.Min; x <= pd.ConstraintX.Max; x++ {
					isActive, exists := pd.Cubes[Point{x, y, z, w}]
					var char string
					if exists && isActive {
						char = ACTIVE_REPR
					} else {
						char = INACTIVE_REPR
					}
					row += char
				}
				fmt.Println(row)
			}
		}
	}
}

func (pd *PocketDimension) GetActiveCubeCount() int {
	count := 0
	for _, isActive := range pd.Cubes {
		if isActive {
			count++
		}
	}

	return count
}

func GetNewState(isCurrentlyActive bool, neighborStates map[bool]int) bool {
	if isCurrentlyActive {
		return neighborStates[ACTIVE] == 3 || neighborStates[ACTIVE] == 2
	} else {
		return neighborStates[ACTIVE] == 3
	}
}

func main() {
	dat, err := ioutil.ReadFile("aoc17.txt")
	if err != nil {
		panic(err)
	}
	txt := string(dat)

	pd := NewPocketDimension(txt)
	pd.Mode = MODE_PART_ONE
	for pd.Cycle < 6 {
		pd.ExecuteCycle()
	}

	fmt.Println("Part one:", pd.GetActiveCubeCount())

	pd = NewPocketDimension(txt)
	pd.Mode = MODE_PART_TWO
	for pd.Cycle < 6 {
		pd.ExecuteCycle()
	}

	fmt.Println("Part two:", pd.GetActiveCubeCount())
}
