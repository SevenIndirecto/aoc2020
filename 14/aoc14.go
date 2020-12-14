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
	Zero  rune = '0'
	One   rune = '1'
	Float rune = 'X'
)

type MaskBit struct {
	Bit  int
	Type rune
}

type Decoder struct {
	Memory  map[int]int
	Mask    []MaskBit
	Version int
}

func NewDecoder(version int) Decoder {
	return Decoder{Memory: make(map[int]int), Version: version}
}

func (decoder *Decoder) ExecuteLine(line string) {
	if len(line) < 2 {
		return
	}

	if line[1] == 'e' {
		decoder.LoadValue(line)
	} else {
		decoder.LoadMask(line)
	}
}

func (decoder *Decoder) LoadMask(instruction string) {
	var maskBits []MaskBit
	s := strings.Split(instruction, " = ")

	mask := []rune(s[1])
	size := len(mask)
	for i := 0; i < size; i++ {
		maskBits = append(maskBits, MaskBit{
			Bit:  i,
			Type: mask[size-1-i],
		})
	}
	decoder.Mask = maskBits
}

func (decoder *Decoder) LoadValue(instruction string) {
	s := strings.Split(instruction, " = ")
	re := regexp.MustCompile(`^mem\[(\d+)]$`)
	match := re.FindStringSubmatch(s[0])
	if match == nil {
		panic(fmt.Sprintf("Invalid mem address [%s]", s[0]))
	}
	address, _ := strconv.Atoi(match[1])
	value, _ := strconv.Atoi(s[1])

	if decoder.Version == 1 {
		decoder.Memory[address] = decoder.applyMaskToValue(value)
	} else {
		decoder.applyMaskToAddress(0, 0, address, value)
	}
}

func (decoder *Decoder) applyMaskToValue(value int) int {
	for _, rule := range decoder.Mask {
		switch rule.Type {
		case Zero:
			value &^= 1 << rule.Bit // clear bit
		case One:
			value |= 1 << rule.Bit // set bit
		}
	}
	return value
}

func (decoder *Decoder) applyMaskToAddress(offset, decodedAddress, address, value int) {
	if offset == len(decoder.Mask) {
		decoder.Memory[decodedAddress] = value
		return
	}

	var bitSet int
	offsetValue := int(math.Pow(2, float64(offset)))

	switch decoder.Mask[offset].Type {
	case Float:
		// Split
		decoder.applyMaskToAddress(offset+1, decodedAddress+offsetValue, address, value)
		decoder.applyMaskToAddress(offset+1, decodedAddress, address, value)
		return
	case Zero:
		// Bit remains unchanged
		bitSet = (address >> offset) & 1
	case One:
		// Set to one
		bitSet = 1
	}
	decoder.applyMaskToAddress(offset+1, decodedAddress+bitSet*offsetValue, address, value)
}

func (decoder *Decoder) Sum() int {
	sum := 0
	for _, val := range decoder.Memory {
		sum += val
	}
	return sum
}

func main() {
	dat, err := ioutil.ReadFile("aoc14.txt")
	if err != nil {
		panic(err)
	}
	txt := string(dat)
	lines := strings.Split(txt, "\n")

	// Part one
	decoder := NewDecoder(1)
	for _, line := range lines {
		decoder.ExecuteLine(line)
	}
	fmt.Println("Part one:", decoder.Sum())

	// Part two
	decoder = NewDecoder(2)
	for _, line := range lines {
		decoder.ExecuteLine(line)
	}
	fmt.Println("Part two:", decoder.Sum())
}
