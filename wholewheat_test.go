package main

import (
	"io/ioutil"
	"log"
	"math/big"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func readRecipe(b []byte) *Recipe {
	var recipe Recipe
	err := yaml.Unmarshal(b, &recipe)
	if err != nil {
		log.Fatal(err)
	}
	return &recipe
}

func readWholeWheatRecipe() *Recipe {
	b, err := ioutil.ReadFile("testdata/whole-wheat-rustic-italian-bread.yml")
	if err != nil {
		log.Fatal(err)
	}
	return readRecipe(b)
}

func measure(t *testing.T, m Measurement, count int64, frac *big.Rat) Measurement {
	mul := big.NewRat(count, 1)
	if frac != nil {
		mul.Add(mul, frac)
	}
	result, err := m.Mul(mul)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestParseFile(t *testing.T) {
	readWholeWheatRecipe()
}

func TestShoppingList(t *testing.T) {
	var expected = []Ingredient{
		{"bread flour", []Measurement{
			measure(t, Ounce, 20, big.NewRat(1, 2)),
			measure(t, Cup, 3, big.NewRat(3, 4))}},
		{"instant yeast", []Measurement{
			measure(t, Teaspoon, 1, big.NewRat(1, 4))}},
		{"water room temp", []Measurement{
			measure(t, Ounce, 18, big.NewRat(7, 10)),
			measure(t, Cup, 2, big.NewRat(1, 3))}},
		{"whole wheat flour", []Measurement{
			measure(t, Ounce, 7, nil),
			measure(t, Cup, 1, big.NewRat(1, 4))}},
		{"table salt", []Measurement{
			measure(t, Teaspoon, 2, nil)}},
	}
	recipe := readWholeWheatRecipe()
	found, err := recipe.ShoppingList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("got %v; want %v", found, expected)
	}
}

func verifyShoppingListCombine(t *testing.T, b []byte, expected []Ingredient) {
	recipe := readRecipe(b)
	found, err := recipe.ShoppingList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("got %v; want %v", found, expected)
	}
}

func TestShoppingListCombine(t *testing.T) {
	verifyShoppingListCombine(t, []byte(combinable),
		[]Ingredient{
			{"bread flour", []Measurement{
				measure(t, Ounce, 20, big.NewRat(1, 2)),
				measure(t, Cup, 3, big.NewRat(3, 4))}}})
	verifyShoppingListCombine(t, []byte(notCombinable),
		[]Ingredient{
			{"bread flour", []Measurement{
				measure(t, Ounce, 11, nil)}},
			{"bread flour", []Measurement{
				measure(t, Cup, 1, big.NewRat(3, 4))}}})
}

const combinable = `
sections:
-
 name: biga
 ingredients:
   - [ bread flour, 11 oz, 2 cups]
-
  ingredients:
    - [ bread flour, 9 1/2 oz, 1 3/4 cups]
`

const notCombinable = `
sections:
-
 name: biga
 ingredients:
   - [ bread flour, 11 oz]
-
  ingredients:
    - [ bread flour, 1 3/4 cups]
`
