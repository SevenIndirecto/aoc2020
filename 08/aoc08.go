package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Instruction struct {
	Op  string
	Arg int
}

type Operation func(arg int)

type GameConsole struct {
	Acc          int
	Ip           int
	Instructions []Instruction
	Executed     map[int]bool
	Ops          map[string]Operation
}

func (gc *GameConsole) Init(instructions []Instruction) {
	gc.Instructions = instructions
	gc.Acc = 0
	gc.Ip = 0
	gc.Executed = make(map[int]bool)
	gc.Ops = map[string]Operation{
		"jmp": func(arg int) {
			gc.UpdateIp(arg)
			gc.Ip--
		},
		"nop": func(arg int) {},
		"acc": func(arg int) {
			gc.Acc += arg
		},
	}
}

func (gc *GameConsole) UpdateIp(delta int) {
	gc.Ip += delta
}

func (gc *GameConsole) Execute() {
	if gc.Ip < 0 || gc.Ip >= len(gc.Instructions) {
		panic("IP overflow error")
	}
	gc.Executed[gc.Ip] = true

	instruction := gc.Instructions[gc.Ip]
	gc.Ops[instruction.Op](instruction.Arg)
	gc.Ip++
}

func (gc *GameConsole) Run() bool {
	for {
		if _, alreadyExecuted := gc.Executed[gc.Ip]; alreadyExecuted {
			return false
		}
		if gc.Ip == len(gc.Instructions) {
			return true
		}
		gc.Execute()
	}
}

func LoadInstructions(path string) []Instruction {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	var instructions []Instruction

	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		ins := strings.Split(line, " ")
		arg, _ := strconv.Atoi(ins[1])
		instructions = append(instructions, Instruction{
			Op:  ins[0],
			Arg: arg,
		})
	}
	return instructions
}

func FixInstructions(instructions []Instruction) int {
	gc := GameConsole{}
	togglePos := -1

	for found := false; !found; {
		togglePos++
		modifiedInstructions := make([]Instruction, len(instructions))
		copy(modifiedInstructions, instructions)

		switch modifiedInstructions[togglePos].Op {
		case "jmp":
			modifiedInstructions[togglePos].Op = "nop"
		case "nop":
			modifiedInstructions[togglePos].Op = "jmp"
		default:
			continue
		}

		gc.Init(modifiedInstructions)
		found = gc.Run()
	}
	return gc.Acc
}

func main() {
	instructions := LoadInstructions("aoc08.txt")
	gc := GameConsole{}
	gc.Init(instructions)
	gc.Run()

	fmt.Println("Part one: ", gc.Acc)
	fmt.Println("Part two: ", FixInstructions(instructions))
}
