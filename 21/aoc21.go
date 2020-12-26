package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type Ingredient string
type Allergen string
type Food struct {
	Allergens   map[Allergen]bool
	Ingredients map[Ingredient]bool
}

func (f *Food) ContainsIngredient(ing Ingredient) bool {
	_, exists := f.Ingredients[ing]
	return exists
}

func (f *Food) ContainsAllergen(alg Allergen) bool {
	_, exists := f.Allergens[alg]
	return exists
}

func (f Food) String() string {
	str := ""
	for ing := range f.Ingredients {
		str += string(ing) + " "
	}
	str += "(contains "
	for alg := range f.Allergens {
		str += string(alg) + ", "
	}
	str = str[:len(str)-2] + ")"
	return str
}

func PrintFoodList(foodList []Food) {
	for _, f := range foodList {
		fmt.Println(f)
	}
}

func GetIngredientsThatCannotContainAllergens(foods *[]Food) []Ingredient {
	allergens := make(map[Allergen]bool)
	ingredients := make(map[Ingredient]bool)

	for _, f := range *foods {
		for ing := range f.Ingredients {
			ingredients[ing] = true
		}
		for alg := range f.Allergens {
			allergens[alg] = true
		}
	}

	var unassignable []Ingredient
	solved := make(map[Ingredient]bool) // should swap unassignable usage with solved, but too lazy to refactor

	for ing := range ingredients {
		assignable := false
		var unassignedIngredients []Ingredient
		for uIng := range ingredients {
			if uIng != ing {
				unassignedIngredients = append(unassignedIngredients, uIng)
			}
		}
		for alg := range allergens {
			var unassignedAllergens []Allergen
			for uAlg := range allergens {
				if uAlg != alg {
					unassignedAllergens = append(unassignedAllergens, uAlg)
				}
			}

			if IsSolvable(unassignedAllergens, unassignedIngredients, map[Allergen]Ingredient{alg: ing}, foods, solved) {
				assignable = true
				break
			}
		}
		if !assignable {
			unassignable = append(unassignable, ing)
			solved[ing] = true
		}
	}

	return unassignable
}

func IsSolvable(
	unassignedAllergens []Allergen,
	unassignedIngredients []Ingredient,
	alg2ing map[Allergen]Ingredient,
	foods *[]Food,
	solved map[Ingredient]bool,
) bool {
	// Validate rules
	if !CanBeAppliedToFoodList(*foods, alg2ing) {
		return false
	}
	if len(unassignedAllergens) < 1 {
		return true
	}
	if len(unassignedIngredients) < 1 {
		return false
	}

	for indexAlg, alg := range unassignedAllergens {
		for indexIng, ing := range unassignedIngredients {
			if _, exists := solved[ing]; exists {
				// Already established this ingredient can't be assigned an allergen
				continue
			}

			newAlg2ing := make(map[Allergen]Ingredient)
			for k, v := range alg2ing {
				newAlg2ing[k] = v
			}
			newAlg2ing[alg] = ing

			algCopy := make([]Allergen, len(unassignedAllergens))
			copy(algCopy, unassignedAllergens)
			copy(algCopy[indexAlg:], algCopy[indexAlg+1:])

			ingCopy := make([]Ingredient, len(unassignedIngredients))
			copy(ingCopy, unassignedIngredients)
			copy(ingCopy[indexIng:], ingCopy[indexIng+1:])

			if IsSolvable(algCopy[:len(algCopy)-1], ingCopy[:len(ingCopy)-1], newAlg2ing, foods, solved) {
				// We found an allergen + ingredient combination reduction that conforms with
				// our food list, so the state supplied to this function is solvable
				return true
			}
		}
	}

	return false
}

// Each allergen is found in exactly one ingredient. An ingredient contains 0 or 1 allergens.
func CanBeAppliedToFoodList(foodList []Food, alg2ing map[Allergen]Ingredient) bool {
	for _, f := range foodList {
		for algInFood := range f.Allergens {
			// If allergen from food list is defined in current "allergen-to-ingredient-map"
			// then the ingredient MUST be part of food
			if ing, algIsMapped := alg2ing[algInFood]; algIsMapped {
				if !f.ContainsIngredient(ing) {
					return false
				}
			}
		}
	}
	return true
}

func LoadFoodList(input string) []Food {
	lines := strings.Split(input, "\n")

	var foodList []Food

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		s := strings.Split(line, " (contains ")
		ingStrings := strings.Split(s[0], " ")
		algStrings := strings.Split(s[1][:len(s[1])-1], ", ")

		f := Food{
			Allergens:   make(map[Allergen]bool),
			Ingredients: make(map[Ingredient]bool),
		}
		for _, ing := range ingStrings {
			f.Ingredients[Ingredient(ing)] = true
		}
		for _, alg := range algStrings {
			f.Allergens[Allergen(alg)] = true
		}
		foodList = append(foodList, f)
	}
	return foodList
}

func CountAppearance(foods *[]Food, ingredients []Ingredient) int {
	count := 0
	for _, f := range *foods {
		for _, ing := range ingredients {
			if _, exists := f.Ingredients[ing]; exists {
				count++
			}
		}
	}
	return count
}

func RemoveIngredientsFromFoodList(foodList *[]Food, ingsToRemove []Ingredient) {
	for _, f := range *foodList {
		f.RemoveIngredients(ingsToRemove)
	}
}

func (f *Food) RemoveIngredients(ingsToRemove []Ingredient) {
	for _, ing := range ingsToRemove {
		delete(f.Ingredients, ing)
	}
}

func Match(
	foods *[]Food,
	allergensToMatch map[Allergen]bool,
	ingredientsToMatch map[Ingredient]bool,
	alg2ing map[Allergen]Ingredient,
) map[Allergen]Ingredient {
	if len(ingredientsToMatch) == 0 && len(allergensToMatch) == 0 {
		// Done
		return alg2ing
	}

	// Take first ingredient and see which allergen matches
	var ingToMatch Ingredient
	for ing := range ingredientsToMatch {
		ingToMatch = ing
		break
	}

	for alg := range allergensToMatch {
		newAlg2ing := make(map[Allergen]Ingredient)
		for k, v := range alg2ing {
			newAlg2ing[k] = v
		}
		newAlg2ing[alg] = ingToMatch

		if !CanBeAppliedToFoodList(*foods, newAlg2ing) {
			continue
		}

		algCopy := make(map[Allergen]bool)
		for k := range allergensToMatch {
			if alg != k {
				algCopy[k] = true
			}
		}
		ingCopy := make(map[Ingredient]bool)
		for k := range ingredientsToMatch {
			if ingToMatch != k {
				ingCopy[k] = true
			}
		}

		// Got a candidate, can we finish?
		matchedAlg2Ing := Match(foods, algCopy, ingCopy, newAlg2ing)
		if len(matchedAlg2Ing) == len(alg2ing) + len(allergensToMatch) {
			return matchedAlg2Ing
		}
	}
	// No match
	return make(map[Allergen]Ingredient)
}

type AlgIngPair struct {
	Alg Allergen
	Ing Ingredient
}

func GetCanonicalList(list map[Allergen]Ingredient) string {
	var pairs []AlgIngPair
	for alg, ing := range list {
		pairs = append(pairs, AlgIngPair{alg, ing})
	}
	// Sort
	sort.SliceStable(pairs, func(i, j int) bool {
		return string(pairs[i].Alg) < string(pairs[j].Alg)
	})
	str := ""
	for _, p := range pairs {
		str += string(p.Ing) + ","
	}
	return str[:len(str)-1]
}

func main() {
	dat, err := ioutil.ReadFile("aoc21.txt")
	if err != nil {
		panic(err)
	}
	txt := string(dat)

	foodList := LoadFoodList(txt)
	unassignable := GetIngredientsThatCannotContainAllergens(&foodList)
	count := CountAppearance(&foodList, unassignable)
	fmt.Println("Part one:", count)

	RemoveIngredientsFromFoodList(&foodList, unassignable)
	fmt.Println()
	PrintFoodList(foodList)

	allergens := make(map[Allergen]bool)
	ingredients := make(map[Ingredient]bool)
	for _, f := range foodList {
		for alg := range f.Allergens {
			allergens[alg] = true
		}
		for ing := range f.Ingredients {
			ingredients[ing] = true
		}
	}
	fmt.Println(allergens, len(allergens))
	fmt.Println(ingredients, len(ingredients))
	solutionSlice := Match(&foodList, allergens, ingredients, make(map[Allergen]Ingredient))
	solution := GetCanonicalList(solutionSlice)

	fmt.Println("\nPart two: ", solution)
}
