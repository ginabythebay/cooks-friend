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

func TestOutput(t *testing.T) {
	type testCase struct {
		expected string
		sys      System
		m        Measurement
	}

	cases := []*testCase{
		&testCase{"1/8 tsp", Imperial, EighthTeaspoon},
		&testCase{"1 tsp", Imperial, Teaspoon},
		&testCase{"5 ml", Metric, Teaspoon},
		&testCase{"1 oz", Imperial, Ounce},
		&testCase{"2 T", Imperial, Volume(2 * Tablespoon)},
		&testCase{"1 cup, 1 tsp", Imperial, Volume(Teaspoon + Cup)},
		&testCase{"3 ml", Metric, Volume(Milliliter * 3)},
		&testCase{"451 ml", Metric, Volume(Milliliter * 451)},
		&testCase{"1.451 l", Metric, Volume(Milliliter * 1451)},
		&testCase{"3/4 c", Imperial, ThreeQuarterCup},
		&testCase{"gal, 1 c", Imperial, Volume(Gallon + Cup)},
		&testCase{"0.5 oz", Imperial, Weight(Ounce / 2)},
		&testCase{"3 lb,y 8 oz", Imperial, Weight(Ounce*16*3 + Ounce*8)},
	}
	for i, c := range cases {
		s := c.m.Output(c.sys)
		if s != c.expected {
			t.Errorf("Expected %v but got %v when converting %#v using %v in test case %v", c.expected, s, c.m, c.sys, i)
		}
	}
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
