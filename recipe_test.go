package main

import "testing"

func TestMerge(t *testing.T) {
	i := Ingredient{"water", []Measurement{Teaspoon}}
	if err := i.Merge(&Ingredient{"flour", []Measurement{Teaspoon}}); err == nil {
		t.Errorf("Expected error merging water and flour and did not get it")
	}
	if err := i.Merge(&Ingredient{"water", []Measurement{Ounce}}); err == nil {
		t.Errorf("Expected error merging teaspoon and ounce did not get it")
	}
	if err := i.Merge(&Ingredient{"water", []Measurement{Teaspoon}}); err != nil {
		t.Fatal(err)
	}
	if i.Measurements[0] != Volume(Teaspoon*2) {
		t.Errorf("Expected %v but got %v", Volume(Teaspoon*2), i.Measurements[0])
	}

}
