package main

import (
	"testing"
)

type Fixture struct {
	Path string
	Expected int
}

type InstructionFixture struct {
	Instruction Instruction
	InitialIp int
	ExpectedIp int
	ExpectedAcc int
}

func TestInstructions(t *testing.T) {
	fixtures := []InstructionFixture{
		{Instruction{"jmp", 4}, 2, 6, 0},
		{Instruction{"jmp", -3}, 4, 1, 0},
		{Instruction{"acc", 3}, 4, 5, 3},
		{Instruction{"acc", -1}, 1, 2, -1},
	}

	for _, fixture := range fixtures {
		gc := GameConsole{}
		var ins []Instruction
		for i := 0; i < fixture.InitialIp; i++ {
			ins = append(ins, Instruction{"nop", 0})
		}
		ins = append(ins, fixture.Instruction)
		gc.Init(ins)
		gc.Ip = fixture.InitialIp
		gc.Execute()

		if gc.Ip != fixture.ExpectedIp || gc.Acc != fixture.ExpectedAcc {
			t.Errorf(
				"Op: %v = ip: %d acc: %d; want ip: %d acc: %d",
				fixture.Instruction,
				gc.Ip,
				gc.Acc,
				fixture.ExpectedIp,
				fixture.ExpectedAcc,
			)
		}
	}
}

func TestStopAtDuplicate(t *testing.T) {
	fixtures := []Fixture{
		{"aoc08_test1.txt", 5},
	}

	for _, fixture := range fixtures {
		ins := LoadInstructions(fixture.Path)
		gc := GameConsole{}
		gc.Init(ins)
		gc.Run()
		got := gc.Acc

		if got != fixture.Expected {
			t.Errorf("Fixture[%s] = %d; want %d", fixture.Path, got, fixture.Expected)
		}
	}
}

func TestFixInstructions(t *testing.T) {
	fixtures := []Fixture{
		{"aoc08_test1.txt", 8},
	}

	for _, fixture := range fixtures {
		ins := LoadInstructions(fixture.Path)
		got := FixInstructions(ins)

		if got != fixture.Expected {
			t.Errorf("Fixture[%s] = %d; want %d", fixture.Path, got, fixture.Expected)
		}
	}
}
