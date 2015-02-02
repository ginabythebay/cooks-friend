package main

import (
	"math/big"
	"testing"
)

func TestParseMagnitude(t *testing.T) {
	verifyParseMagnitude(t, "1/2", big.NewRat(1, 2))
	verifyParseMagnitude(t, "1 1/2", big.NewRat(3, 2))
	verifyParseMagnitude(t, "2.5", big.NewRat(5, 2))
	verifyParseMagnitude(t, "3", big.NewRat(3, 1))
}

func TestMul(t *testing.T) {
	if v, err := mul(20, big.NewRat(3, 4)); err != nil {
		t.Fatal(err)
	} else {
		if v != 15 {
			t.Errorf("Error multiplying 20 and 3/4: expected 15 and got %v.", v)
		}
	}
	if _, err := mul(20, big.NewRat(1, 7)); err == nil {
		t.Errorf("Expected error muliplying 20 and 1/7 and did not get it.")
	}
}

func TestParse(t *testing.T) {
	verifyParse(t, "1/2 tsp", HalfTeaspoon)
	verifyParse(t, "1 1/2 tsp", Volume(Teaspoon+HalfTeaspoon))
	verifyParse(t, "3/4 gallons", Volume(Quart*3))
	verifyParse(t, "16 oz", Pound)
	verifyParse(t, "1/2 lb", Weight(Ounce*8))
}

func verifyParseMagnitude(t *testing.T, s string, expected *big.Rat) {
	if value, err := parseMagnitude(s); err != nil {
		t.Fatal(err)
	} else {
		if value.Num().Int64() != expected.Num().Int64() ||
			value.Denom().Int64() != expected.Denom().Int64() {
			t.Errorf("Error parsing magnitude %v: expected %v but got %v", s, expected, value)
		}
	}
}

func verifyParse(t *testing.T, s string, expected Measurement) {
	if value, err := Parse(s); err != nil {
		t.Fatal(err)
	} else {
		if value != expected {
			t.Errorf("Error parsing %v: expected %v but got %v", s, expected, value)
		}
	}
}
