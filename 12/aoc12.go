package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

const (
	North = iota
	East
	South
	West
)

type Point struct {
	X, Y int
}

type Instruction struct {
	Action string
	Value int
}

type Ship struct {
	Position Point
	Direction int
	Instructions []Instruction
	Waypoint Point
	Ip int
}

func (ship *Ship) MoveForward(distance int) {
	switch ship.Direction {
	case North:
		ship.Position.Y += distance
	case East:
		ship.Position.X += distance
	case South:
		ship.Position.Y -= distance
	case West:
		ship.Position.X -= distance
	}
}

func (ship *Ship) Turn(instruction Instruction) {
	turns := instruction.Value / 90

	clockwiseModifier := 1
	if instruction.Action == "L" {
		clockwiseModifier = -1
	}

	// NOTE: % returns negative numbers, might write a helper if this ends up repeating
	ship.Direction = (((ship.Direction + clockwiseModifier * turns) % 4) + 4) % 4
}

func (ship *Ship) movePoint(point *Point, instruction Instruction) {
	switch instruction.Action {
	case "N":
		point.Y += instruction.Value
	case "E":
		point.X += instruction.Value
	case "S":
		point.Y -= instruction.Value
	case "W":
		point.X -= instruction.Value
	}
}

func (ship *Ship) MoveInDirection(instruction Instruction) {
	ship.movePoint(&ship.Position, instruction)
}

func (ship *Ship) Execute() {
	ins := ship.Instructions[ship.Ip]
	switch ins.Action {
	case "N", "E", "S", "W":
		ship.MoveInDirection(ins)
	case "R", "L":
		ship.Turn(ins)
	case "F":
		ship.MoveForward(ins.Value)
	}
	//fmt.Println(ins, ship)
}

func (ship *Ship) ExecuteReal() {
	ins := ship.Instructions[ship.Ip]
	switch ins.Action {
	case "N", "E", "S", "W":
		ship.MoveWaypoint(ins)
	case "R", "L":
		ship.RotateWaypoint(ins)
	case "F":
		ship.MoveToWaypoint(ins)
	}
}

func (ship *Ship) MoveWaypoint(instruction Instruction) {
	ship.movePoint(&ship.Waypoint, instruction)
}

func (ship *Ship) RotateWaypoint(instruction Instruction) {
	times := float64(instruction.Value / 90)
	angle := times * math.Pi / 2
	if instruction.Action == "R" {
		angle = -1 * angle
	}

	ship.Waypoint = Point{
		X: ship.Waypoint.X * int(math.Cos(angle)) - ship.Waypoint.Y * int(math.Sin(angle)),
		Y: ship.Waypoint.X * int(math.Sin(angle)) + ship.Waypoint.Y * int(math.Cos(angle)),
	}
}

func (ship *Ship) MoveToWaypoint(instruction Instruction) {
	ship.Position.X += ship.Waypoint.X * instruction.Value
	ship.Position.Y += ship.Waypoint.Y * instruction.Value
}

func (ship *Ship) DistanceTravelled() int {
	return int(math.Abs(float64(ship.Position.X))) + int(math.Abs(float64(ship.Position.Y)))
}

func (ship *Ship) LoadInstructionSet(lines []string) {
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		value, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(fmt.Sprintf("Invalid line: [%s]", line))
		}
		ins := Instruction{
			Action: line[0:1],
			Value: value,
		}
		ship.Instructions = append(ship.Instructions, ins)
	}
}

func (ship *Ship) ExecuteInstructions() {
	for ship.Ip < len(ship.Instructions) {
		ship.Execute()
		ship.Ip++
	}
}

func (ship *Ship) ExecuteRealInstructions() {
	for ship.Ip < len(ship.Instructions) {
		ship.ExecuteReal()
		ship.Ip++
	}
}

func main() {
	dat, err := ioutil.ReadFile("aoc12.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	ship := Ship{Direction: East, Position: Point{0, 0}}
	ship.LoadInstructionSet(lines)
	ship.ExecuteInstructions()
	fmt.Println("Part one:", ship.DistanceTravelled())

	ship = Ship{Waypoint: Point{10, 1}, Position: Point{0, 0}}
	ship.LoadInstructionSet(lines)
	ship.ExecuteRealInstructions()
	fmt.Println("Part two:", ship.DistanceTravelled())
}
