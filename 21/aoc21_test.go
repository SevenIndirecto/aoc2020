package main

import (
	"reflect"
	"testing"
)

type Fixture struct {
	Alg2ing map[Allergen]Ingredient
	Expected bool
}

const FOODS = `mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
trh fvjkl sbzzf mxmxvkd (contains dairy)
sqjhc fvjkl (contains soy)
sqjhc mxmxvkd sbzzf (contains fish)`

func TestCanBeAppliedToFoodList(t *testing.T) {
	foodList := LoadFoodList(FOODS)
	fixtures := []Fixture {
		{
			map[Allergen]Ingredient{
				"soy": "fvjkl",
				"dairy": "mxmxvkd",
				"fish": "sqjhc",
			},
			true,
		},
		{
			map[Allergen]Ingredient{
				"soy": "sqjhc",
				"dairy": "fvjkl",
				"fish": "mxmxvkd",
			},
			false,
		},
		{
			map[Allergen]Ingredient{
				"fish": "fvjkl",
				"soy": "sbzzf",
			},
			false,
		},
		{
			map[Allergen]Ingredient{
				"fish": "fvjkl",
				"dairy": "sqjhc",
			},
			false,
		},
	}

	for _, f := range fixtures {
		got := CanBeAppliedToFoodList(foodList, f.Alg2ing)
		if got != f.Expected {
			t.Errorf("Trying %v got %v expected %v", f.Alg2ing, got, f.Expected)
		}
	}
}

func TestGetIngredientsThatCannotContainAllergens(t *testing.T) {
	foodList := LoadFoodList(FOODS)
	got := GetIngredientsThatCannotContainAllergens(&foodList)
	expected := map[Ingredient]bool{"kfcds": true, "nhms": true, "trh": true, "sbzzf": true}

	if len(got) != len(expected) {
		t.Errorf("Got unexpected ingredients %v expected %v", got, expected)
	}
	for _, ing := range got {
		if _, exists := expected[ing]; !exists {
			t.Errorf("Unexpected ingredient %v", ing)
		}
	}
}

func TestCountAppearance(t *testing.T) {
	foodList := LoadFoodList(FOODS)
	ingredients := []Ingredient{"kfcds", "nhms", "trh", "sbzzf"}
	got := CountAppearance(&foodList, ingredients)
	expected := 5

	if got != expected {
		t.Errorf("Got %d expected %d", got , expected)
	}
}

func TestRemoveIngredientsFromFoodList(t *testing.T) {
	foodList := LoadFoodList(FOODS)
	ingsToRemove := GetIngredientsThatCannotContainAllergens(&foodList)
	RemoveIngredientsFromFoodList(&foodList, ingsToRemove)
	expected := []Food{
		{
			map[Allergen]bool{"fish": true, "dairy": true},
			map[Ingredient]bool{"mxmxvkd": true, "sqjhc": true},
		},
		{
			map[Allergen]bool{"dairy": true},
			map[Ingredient]bool{"fvjkl": true, "mxmxvkd": true},
		},
		{
			map[Allergen]bool{"soy": true},
			map[Ingredient]bool{"fvjkl": true, "sqjhc": true},
		},
		{
			map[Allergen]bool{"fish": true},
			map[Ingredient]bool{"mxmxvkd": true, "sqjhc": true},
		},
	}

	if !reflect.DeepEqual(foodList, expected) {
		t.Errorf("Failed to remove ingredients")
	}
}

func TestMatch(t *testing.T) {
	foodList := LoadFoodList(FOODS)
	ingsToRemove := GetIngredientsThatCannotContainAllergens(&foodList)
	RemoveIngredientsFromFoodList(&foodList, ingsToRemove)

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
	gotSlice := Match(&foodList, allergens, ingredients, make(map[Allergen]Ingredient))
	got := GetCanonicalList(gotSlice)
	expected := "mxmxvkd,sqjhc,fvjkl"

	if got != expected {
		t.Errorf("Got %s expected %s", got, expected)
	}
}
