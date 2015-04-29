package main

import (
	"io/ioutil"
	"log"
	"math/big"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func readRecipe() *Recipe {
	b, err := ioutil.ReadFile("testdata/whole-wheat-rustic-italian-bread.yml")
	if err != nil {
		log.Fatal(err)
	}
	var recipe Recipe
	err = yaml.Unmarshal(b, &recipe)
	if err != nil {
		log.Fatal(err)
	}
	return &recipe
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
	readRecipe()
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
	recipe := readRecipe()
	found, err := recipe.ShoppingList()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("got %v; want %v", found, expected)
	}
}
