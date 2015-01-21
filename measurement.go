package main

import (
	"fmt"
	"log"
)
import "regexp"
import "reflect"

var (
	re = regexp.MustCompile(`^\s*(\d+\s+\d+/\d+|\d+|\d+\.\d+|\d+/\d+)\s+([a-zA-Z]*)\s*$`)
)

type Volume int64

const (
	Milliliter Volume = 8
	Deciliter         = Milliliter * 100
	Liter             = Milliliter * 1000

	EighthTeaspoon  = Milliliter * 5 / 8
	QuarterTeaspoon = EighthTeaspoon * 2
	HalfTeaspoon    = QuarterTeaspoon * 2
	Teaspoon        = HalfTeaspoon * 2
	Tablespoon      = Teaspoon * 3
	FluidOunce      = Tablespoon * 2
	Cup             = FluidOunce * 8
	Pint            = Cup * 2
	Quart           = Cup * 4
	Gallon          = Quart * 4
)

func (v Volume) Add(o Measurement) (result Measurement, err error) {
	if other, ok := o.(Volume); ok {
		return v + other, nil
	} else {
		return v, fmt.Errorf("Volume incompatible with %v", reflect.ValueOf(o).Type())
	}
}

type Weight int64

const (
	Milligram Weight = 8
	Gram             = Milligram * 1000
	Kilogram         = Gram * 1000

	Ounce = 28409 * Milligram
	Pound = Ounce * 16
)

func (v Weight) Add(o Measurement) (result Measurement, err error) {
	if other, ok := o.(Weight); ok {
		return v + other, nil
	} else {
		return v, fmt.Errorf("Weight incompatible with %v", reflect.ValueOf(o).Type())
	}
}

type Measurement interface {
	Add(other Measurement) (result Measurement, err error)
}

func Parse(s string) (m Measurement, err error) {
	if matches := re.FindStringSubmatch(s); matches != nil {
		if len(matches) == 3 {
			magnitude := matches[1]
			units := matches[2]
			log.Printf("%#v: %#v: %#v", s, magnitude, units)
			return nil, nil
		} else {
			return nil, fmt.Errorf("Unable to parse [%v] as measurment.  Matches was %#v", s, matches)
		}
	} else {
		return nil, fmt.Errorf("Unable to parse [%v] as measurement via regexp", s)
	}

}

// need maps of accepted strings to the matching values

// Is there some kind of 'unit' interface that these both implement?
// Seems like we want a method that accepts a string and returns an
// insance of that unit
