package main

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"regexp"
	"strings"
)

var (
	// we match <magnitude> <unit>
	// where <magnitude> can look like:
	//   2
	//   2.5
	//   2 1/2
	// and <unit> can look like
	//   oz
	//   cup
	//   cups
	//   tsp
	//   etc...
	re = regexp.MustCompile(`^\s*(\d+\s+\d+/\d+|\d+|\d+\.\d+|\d+/\d+)\s+([a-zA-Z]*)\s*$`)
)

type System int8

const (
	Metric System = iota
	Imperial
)

type OutputType int8

const (
	NeverOutput = iota
	SingleOnly
	MultiplesOK
	MultiplesOKAppendS
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
	QuarterCup      = FluidOunce * 2
	HalfCup         = FluidOunce * 4
	Cup             = FluidOunce * 8
	Pint            = Cup * 2
	Quart           = Cup * 4
	Gallon          = Quart * 4
)

type unitInfo struct {
	measurement   Measurement
	system        System
	out           string
	in            []string
	decimalPlaces int8
	outputType    OutputType
}

var (
	volumeInfo = map[Measurement]*unitInfo{
		Milliliter: &unitInfo{
			nil,
			Metric,
			"ml",
			[]string{"ml", "milliliter", "milliliters", "millilitre", "millilitres", "mL"},
			0,
			MultiplesOK,
		},
		Deciliter: &unitInfo{
			nil,
			Metric,
			"",
			[]string{"dl", "deciliter", "deciliters", "decilitre", "decilitres", "dL"},
			0,
			NeverOutput,
		},
		Liter: &unitInfo{
			nil,
			Metric,
			"l",
			[]string{"l", "liter", "liters", "litre", "litres", "L"},
			3,
			MultiplesOK,
		},

		EighthTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/8 tsp",
			[]string{},
			0,
			SingleOnly,
		},
		QuarterTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/4 tsp",
			[]string{},
			0,
			SingleOnly,
		},
		HalfTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/2 tsp",
			[]string{},
			0,
			SingleOnly,
		},
		Teaspoon: &unitInfo{
			nil,
			Imperial,
			"t",
			[]string{"t", "teaspoon", "teaspoons", "tsp.", "tsp"},
			0,
			SingleOnly,
		},
		Tablespoon: &unitInfo{
			nil,
			Imperial,
			"T",
			[]string{"T", "tablespoon", "tablespoons", "tbl.", "tbl", "tbs.", "tbsp."},
			0,
			SingleOnly,
		},
		FluidOunce: &unitInfo{
			nil,
			Imperial,
			"",
			[]string{"fluid ounce", "fluid ounces", "fl oz"},
			0,
			NeverOutput,
		},
		QuarterCup: &unitInfo{
			nil,
			Imperial,
			"1/4 cup",
			[]string{},
			0,
			SingleOnly,
		},
		HalfCup: &unitInfo{
			nil,
			Imperial,
			"1/2 cup",
			[]string{},
			0,
			SingleOnly,
		},
		Cup: &unitInfo{
			nil,
			Imperial,
			"c",
			[]string{"c", "cup", "cups"},
			0,
			MultiplesOK,
		},
		Pint: &unitInfo{
			nil,
			Imperial,
			"",
			[]string{"p", "pt", "pint", "pints", "fl pt"},
			0,
			NeverOutput,
		},
		Quart: &unitInfo{
			nil,
			Imperial,
			"qt",
			[]string{"q", "quart", "quarts", "qt", "fl qt"},
			0,
			MultiplesOK,
		},
		Gallon: &unitInfo{
			nil,
			Imperial,
			"gal",
			[]string{"gal", "gallon", "gallons", "g"},
			0,
			MultiplesOK,
		},
	}
)

func (v Volume) Add(o Measurement) (result Measurement, err error) {
	if other, ok := o.(Volume); ok {
		return v + other, nil
	} else {
		return v, fmt.Errorf("Volume incompatible with %v", reflect.ValueOf(o).Type())
	}
}

func (v Volume) Mul(r *big.Rat) (Measurement, error) {
	if i, err := mul(int64(v), r); err != nil {
		return Teaspoon, err
	} else {
		return Volume(i), nil
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

func (v Weight) Mul(r *big.Rat) (Measurement, error) {
	if i, err := mul(int64(v), r); err != nil {
		return Ounce, err
	} else {
		return Weight(i), nil
	}
}

type Measurement interface {
	Add(other Measurement) (result Measurement, err error)
	Mul(r *big.Rat) (Measurement, error)
}

func mul(i int64, r *big.Rat) (int64, error) {
	result := big.NewRat(i, 1)
	result.Mul(result, r)
	if !result.IsInt() {
		return 0, fmt.Errorf("Error multiplying %v by %v.  We ended up with non-integral value %v", i, r, result)
	}
	return result.Num().Int64(), nil
}

// Can parse things like "1/4", ".5", "2", "2 1/2"
func parseMagnitude(s string) (*big.Rat, error) {
	tokens := strings.Split(s, " ")
	accum := new(big.Rat)
	r := new(big.Rat)
	for _, t := range tokens {
		if _, ok := r.SetString(t); !ok {
			return nil, fmt.Errorf("Error parsing %v.  Token %v could not be parsed by big.Rat")
		}
		accum.Add(accum, r)
	}
	return accum, nil
}

func Parse(s string) (m Measurement, err error) {
	if matches := re.FindStringSubmatch(s); matches != nil {
		if len(matches) == 3 {
			magnitude := matches[1]
			unit := matches[2]
			log.Printf("%#v: %#v: %#v", s, magnitude, unit)
			if info, ok := measurementLookup[unit]; ok {
				if mag, err := parseMagnitude(magnitude); err != nil {
					return nil, err
				} else {
					return info.measurement.Mul(mag)
				}
			} else {
				return nil, fmt.Errorf("Could not recognize [%v] as unit in %s", unit, s)
			}
		} else {
			return nil, fmt.Errorf("Unable to parse [%v] as measurment.  Matches was %#v", s, matches)
		}
	} else {
		return nil, fmt.Errorf("Unable to parse [%v] as measurement via regexp", s)
	}

}

// TODO(gina) add weight here when we have it
var allMeasurementInfo = []map[Measurement]*unitInfo{volumeInfo}

var measurementLookup = make(map[string]*unitInfo)

func init() {
	for _, m := range allMeasurementInfo {
		for measurement, info := range m {
			info.measurement = measurement

			for _, in := range info.in {
				measurementLookup[in] = info
			}
		}
	}
}

// 	}
// }

// need maps of accepted strings to the matching values

// Is there some kind of 'unit' interface that these both implement?
// Seems like we want a method that accepts a string and returns an
// insance of that unit
