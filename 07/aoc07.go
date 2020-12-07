package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type BagCapacity struct {
	Qty   int
	Color string
}

type Bag struct {
	Color                string
	Capacity             []BagCapacity
	KnownToContainTarget bool
}

func ParseBagRule(rule string) Bag {
	s := strings.Split(rule, " contain ")
	re := regexp.MustCompile(`(\w+ \w+) bags`)
	match := re.FindStringSubmatch(s[0])
	color := match[1]

	re = regexp.MustCompile(`(\d+) (\w+ \w+) bag`)
	matches := re.FindAllStringSubmatch(s[1], -1)
	var capacities []BagCapacity

	for i := 0; i < len(matches); i++ {
		qty, _ := strconv.Atoi(matches[i][1])
		c := BagCapacity{
			Qty:   qty,
			Color: matches[i][2],
		}
		capacities = append(capacities, c)
	}

	return Bag{
		Color:                color,
		Capacity:             capacities,
		KnownToContainTarget: false,
	}
}

func CanContainColor(bag Bag, color string, bags map[string]Bag) bool {
	if bag.KnownToContainTarget {
		return true
	}

	for _, bagCapacity := range bag.Capacity {
		containedBag := bags[bagCapacity.Color]
		if containedBag.Color == color || CanContainColor(containedBag, color, bags) {
			bag.KnownToContainTarget = true
			return true
		}
	}
	return false
}

func LoadRules(path string) map[string]Bag {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	bags := make(map[string]Bag)

	for _, rule := range lines {
		if len(rule) < 2 {
			continue
		}
		bag := ParseBagRule(rule)
		bags[bag.Color] = bag
	}
	return bags
}

func GetBagsRequiredForBag(bag Bag, bags map[string]Bag) int {
	requiredChildren := 0
	for _, bagCapacity := range bag.Capacity {
		containedBag := bags[bagCapacity.Color]

		requiredChildren += bagCapacity.Qty
		requiredChildren += bagCapacity.Qty * GetBagsRequiredForBag(containedBag, bags)
	}
	return requiredChildren
}

func main() {
	bags := LoadRules("aoc07.txt")
	count := 0
	for _, bag := range bags {
		if CanContainColor(bag, "shiny gold", bags) {
			count++
		}
	}
	fmt.Println("Part one: ", count)
	targetBag := bags["shiny gold"]
	fmt.Println("Part two: ", GetBagsRequiredForBag(targetBag, bags))
}
