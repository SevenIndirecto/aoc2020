package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Passport map[string]string

func LoadPassports(path string) []Passport {
	var passports []Passport

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	lines := strings.Split(txt, "\n")

	pass := Passport{}

	for _, line := range lines {
		if len(line) < 2 {
			passports = append(passports, pass)
			pass = Passport{}
			continue
		}

		entries := strings.Split(line, " ")
		for _, entry := range entries {
			item := strings.Split(entry, ":")
			pass[item[0]] = item[1]
		}
	}

	return passports
}

func HasRequired(pass Passport, required []string) bool {
	count := 0
	for _, key := range required {
		if _, ok := pass[key]; ok {
			count++
		}
	}
	return len(required) == count
}

type Validator func(string) bool

func IsYearBetween(val string, min int, max int) bool {
	found, _ := regexp.MatchString(`^\d{4}$`, val)
	if !found {
		return false
	}
	n, _ := strconv.Atoi(val)
	return n >= min && n <= max
}

func HasValidValues(pass Passport, validators map[string]Validator) bool {
	for key, value := range pass {
		validate, ok := validators[key]
		if !ok {
			continue
		}
		valid := validate(value)
		if !valid {
			return false
		}
	}
	return true
}

func Validate(passports []Passport) (int, int) {
	required := []string{
		"byr",
		"iyr",
		"eyr",
		"hgt",
		"hcl",
		"ecl",
		"pid",
	}

	validators := map[string]Validator{
		"byr": func(val string) bool {
			return IsYearBetween(val, 1920, 2002)
		},
		"iyr": func(val string) bool {
			return IsYearBetween(val, 2010, 2020)
		},
		"eyr": func(val string) bool {
			return IsYearBetween(val, 2020, 2030)
		},
		"hgt": func(val string) bool {
			re := regexp.MustCompile(`^(\d+)(in|cm)$`)
			match := re.FindStringSubmatch(val)
			if match == nil {
				return false
			}
			h, _ := strconv.Atoi(match[1])
			if match[2] == "cm" {
				return h >= 150 && h <= 193
			}
			if match[2] == "in" {
				return h >= 59 && h <= 76
			}
			return false
		},
		"hcl": func(val string) bool {
			found, _ := regexp.MatchString(`^#[0-9a-f]{6}$`, val)
			return found
		},
		"ecl": func(val string) bool {
			found, _ := regexp.MatchString(`^(amb|blu|brn|gry|grn|hzl|oth)$`, val)
			return found
		},
		"pid": func(val string) bool {
			found, _ := regexp.MatchString(`^\d{9}$`, val)
			return found
		},
	}

	requiredCount := 0
	validCount := 0
	for _, pass := range passports {
		if HasRequired(pass, required) {
			requiredCount++

			if HasValidValues(pass, validators) {
				validCount++
			}
		}
	}
	return requiredCount, validCount
}

func main() {
	passports := LoadPassports("aoc04.txt")
	one, two := Validate(passports)

	fmt.Println("Part one:", one)
	fmt.Println("Part two:", two)
}
