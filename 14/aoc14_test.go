package main

import (
	"reflect"
	"testing"
)

type Fixture struct {
	Mask            string
	ExpectedMaskBit []MaskBit
}

func TestDecoder_LoadMask(t *testing.T) {
	fixtures := []Fixture{
		{
			"mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X",
			[]MaskBit{
				{Bit: 1, Type: Zero},
				{Bit: 6, Type: One},
			},
		},
	}

	for _, f := range fixtures {
		decoder := NewDecoder(1)
		decoder.LoadMask(f.Mask)
		got := decoder.Mask

		for _, maskBit := range f.ExpectedMaskBit {
			if got[maskBit.Bit].Type != maskBit.Type {
				t.Errorf("Parsing(%s) got %v expected %v", f.Mask, got, f.ExpectedMaskBit)
			}
		}
	}
}

type FixtureLoadNumber struct {
	Instruction    string
	ExpectedMemory map[int]int
}

func TestDecoder_LoadNumber(t *testing.T) {
	fixtures := []FixtureLoadNumber{
		{"mem[8] = 11", map[int]int{8: 73}},
		{"mem[7] = 101", map[int]int{7: 101}},
		{"mem[8] = 0", map[int]int{8: 64}},
	}

	for _, f := range fixtures {
		decoder := NewDecoder(1)
		decoder.LoadMask("mask = XXXXXXXXXXXXXXXXXXXXXXXXXXXXX1XXXX0X")
		decoder.ExecuteLine(f.Instruction)
		//decoder.LoadValue(f.Instruction)

		got := decoder.Memory
		if !reflect.DeepEqual(got, f.ExpectedMemory) {
			t.Errorf("LoadNumber(%s) got %v expected %v", f.Instruction, got, f.ExpectedMemory)
		}
	}
}

func TestDecoder_LoadMask2(t *testing.T) {
	lines := []string{
		"mask = 000000000000000000000000000000X1001X",
		"mem[42] = 100",
		"mask = 00000000000000000000000000000000X0XX",
		"mem[26] = 1",
	}
	expectedSum := 208

	decoder := NewDecoder(2)
	for _, line := range lines {
		decoder.ExecuteLine(line)
	}

	got := decoder.Sum()

	if got != expectedSum {
		t.Errorf("Got %d expected %d, state %v", got, expectedSum, decoder)
	}
}
