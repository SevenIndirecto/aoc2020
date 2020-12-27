package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func transform(subjectNumber, loopSize, startAtLoop, startValue int) int {
	value := startValue

	for i := startAtLoop; i < loopSize; i++ {
		value *= subjectNumber
		value %= 20201227
	}
	return value
}

func findLoopSize(pubKey int) int {
	subjectNumber := 7
	loopSize := 0
	transformedValue := 1

	for {
		loopSize++
		transformedValue = transform(subjectNumber, loopSize, loopSize-1, transformedValue)

		if transformedValue == pubKey {
			break
		}
	}
	return loopSize
}

func FindEncryptionKey(doorPubKey, cardPubKey int) int {
	doorLoops := findLoopSize(doorPubKey)
	cardLoops := findLoopSize(cardPubKey)

	encryptionKey := transform(doorPubKey, cardLoops, 0, 1)
	encryptionKeyCheck := transform(cardPubKey, doorLoops, 0, 1)

	if encryptionKey != encryptionKeyCheck {
		panic("Failed to find encryption key")
	}
	return encryptionKey
}

func main() {
	dat, err := ioutil.ReadFile("aoc25.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")
	cardPubKey, _ := strconv.Atoi(lines[0])
	doorPubKey, _ := strconv.Atoi(lines[1])
	fmt.Println("Part one:", FindEncryptionKey(cardPubKey, doorPubKey))
}
