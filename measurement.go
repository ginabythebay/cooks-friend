package main

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"sort"
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

func (sys System) String() string {
	switch sys {
	case Metric:
		return "Metric system"
	case Imperial:
		return "Imperial system"
	}
	return fmt.Sprintf("Unknown system %v", sys)
}

type OutputType byte

// semi-overlapping flags
const (
	SuppressOutput OutputType = 1 << iota // rest ignored if this is set
	MultiplesOk                           // If set, we can have more than one of this unit
	FractionsOk                           // If set, we can have less than one of this unit
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
	ThreeQuarterCup = FluidOunce * 6
	Cup             = FluidOunce * 8
	OneThirdCup     = Cup / 3
	TwoThirdsCups   = Cup * 2 / 3
	Pint            = Cup * 2
	Quart           = Cup * 4
	Gallon          = Quart * 4
)

type unitInfo struct {
	measurement   Measurement
	system        System
	out           string
	in            []string
	decimalPlaces int
	outputType    OutputType
}

type byMeasurement []*unitInfo

func (a byMeasurement) Len() int           { return len(a) }
func (a byMeasurement) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byMeasurement) Less(i, j int) bool { return a[i].measurement.Int64() < a[j].measurement.Int64() }

var (
	volumeInfo = map[Measurement]*unitInfo{
		Milliliter: &unitInfo{
			nil,
			Metric,
			"ml",
			[]string{"ml", "milliliter", "milliliters", "millilitre", "millilitres", "mL"},
			0,
			MultiplesOk,
		},
		Deciliter: &unitInfo{
			nil,
			Metric,
			"",
			[]string{"dl", "deciliter", "deciliters", "decilitre", "decilitres", "dL"},
			0,
			SuppressOutput,
		},
		Liter: &unitInfo{
			nil,
			Metric,
			"l",
			[]string{"l", "liter", "liters", "litre", "litres", "L"},
			3,
			MultiplesOk,
		},

		EighthTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/8 tsp",
			[]string{},
			0,
			0,
		},
		QuarterTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/4 tsp",
			[]string{},
			0,
			0,
		},
		HalfTeaspoon: &unitInfo{
			nil,
			Imperial,
			"1/2 tsp",
			[]string{},
			0,
			0,
		},
		Teaspoon: &unitInfo{
			nil,
			Imperial,
			"tsp",
			[]string{"t", "teaspoon", "teaspoons", "tsp.", "tsp"},
			0,
			MultiplesOk,
		},
		Tablespoon: &unitInfo{
			nil,
			Imperial,
			"T",
			[]string{"T", "tablespoon", "tablespoons", "tbl.", "tbl", "tbs.", "tbsp."},
			0,
			MultiplesOk,
		},
		FluidOunce: &unitInfo{
			nil,
			Imperial,
			"",
			[]string{"fluid ounce", "fluid ounces", "fl oz"},
			0,
			SuppressOutput,
		},
		QuarterCup: &unitInfo{
			nil,
			Imperial,
			"1/4 c",
			[]string{},
			0,
			0,
		},
		HalfCup: &unitInfo{
			nil,
			Imperial,
			"1/2 c",
			[]string{},
			0,
			0,
		},
		ThreeQuarterCup: &unitInfo{
			nil,
			Imperial,
			"3/4 c",
			[]string{},
			0,
			0,
		},
		OneThirdCup: &unitInfo{
			nil,
			Imperial,
			"1/3 c",
			[]string{},
			0,
			0,
		},
		TwoThirdsCups: &unitInfo{
			nil,
			Imperial,
			"2/3 c",
			[]string{},
			0,
			0,
		},
		Cup: &unitInfo{
			nil,
			Imperial,
			"c",
			[]string{"c", "cup", "cups"},
			0,
			MultiplesOk,
		},
		Pint: &unitInfo{
			nil,
			Imperial,
			"",
			[]string{"p", "pt", "pint", "pints", "fl pt"},
			0,
			SuppressOutput,
		},
		Quart: &unitInfo{
			nil,
			Imperial,
			"qt",
			[]string{"q", "quart", "quarts", "qt", "fl qt"},
			0,
			MultiplesOk,
		},
		Gallon: &unitInfo{
			nil,
			Imperial,
			"gal",
			[]string{"gal", "gallon", "gallons", "g"},
			0,
			MultiplesOk,
		},
	}
)

func (v Volume) Output(sys System) string {
	return Output(v, volumeOutput[sys])
}

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

func (v Volume) Int64() int64 {
	return int64(v)
}

type Weight int64

const (
	Milligram Weight = 8 * 10
	Gram             = Milligram * 1000
	Kilogram         = Gram * 1000

	Ounce = 28409 * Milligram
	Pound = Ounce * 16
)

var (
	weightInfo = map[Measurement]*unitInfo{
		Milligram: &unitInfo{
			nil,
			Metric,
			"mg",
			[]string{"mg", "milligram", "milligrams", "milligramme", "milligrammes"},
			0,
			MultiplesOk,
		},
		Gram: &unitInfo{
			nil,
			Metric,
			"g",
			[]string{"g", "gram", "grams", "gramme", "grammes"},
			3,
			MultiplesOk,
		},
		Kilogram: &unitInfo{
			nil,
			Metric,
			"kg",
			[]string{"kg", "kilogram", "kilograms", "kilogramme", "kilogrammes"},
			3,
			MultiplesOk,
		},

		Ounce: &unitInfo{
			nil,
			Imperial,
			"oz",
			[]string{"oz", "ounce", "ounces"},
			1,
			MultiplesOk | FractionsOk,
		},
		Pound: &unitInfo{
			nil,
			Imperial,
			"lb",
			[]string{"lb", "#", "pound", "pounds"},
			0,
			MultiplesOk,
		},
	}
)

func (w Weight) Output(sys System) string {
	return Output(w, weightOutput[sys])
}

func (w Weight) Add(o Measurement) (result Measurement, err error) {
	if other, ok := o.(Weight); ok {
		return w + other, nil
	} else {
		return w, fmt.Errorf("Weight incompatible with %v", reflect.ValueOf(o).Type())
	}
}

func (w Weight) Mul(r *big.Rat) (Measurement, error) {
	if i, err := mul(int64(w), r); err != nil {
		return Ounce, err
	} else {
		return Weight(i), nil
	}
}

func (w Weight) Int64() int64 {
	return int64(w)
}

type Measurement interface {
	Add(other Measurement) (result Measurement, err error)
	Mul(r *big.Rat) (Measurement, error)
	Output(sys System) string
	Int64() int64
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

func Output(m Measurement, units byMeasurement) string {
	remainder := m.Int64()

	// units are already ordered by size, biggest first.  We loop
	// through them, looking for ones to apply, building up a slice of
	// tokens, which will be joined at the end.
	tokens := make([]string, 0)
	for _, u := range units {
		ui := u.measurement.Int64()
		div := remainder / ui
		mod := remainder % ui
		multiplesOK := u.outputType&MultiplesOk != 0
		fractionsOk := u.outputType&FractionsOk != 0
		switch {
		case remainder == ui && !multiplesOK:
			tokens = append(tokens, u.out)
			remainder = 0
		case ui <= remainder && multiplesOK && (mod == 0 || u.decimalPlaces == 0):
			tokens = append(tokens, fmt.Sprintf("%v %s", div, u.out))
			remainder -= div * ui
		case (ui < remainder && u.decimalPlaces != 0) || fractionsOk:
			f := float64(remainder) / float64(ui)
			tokens = append(tokens, fmt.Sprintf("%.*f %s", u.decimalPlaces, f, u.out))
			remainder = 0
		default:
			continue
		}
		if remainder == 0 {
			break
		}
	}

	if len(tokens) != 0 {
		return strings.Join(tokens, ", ")
	} else {
		return "TODO(gina) what to do here?"
	}
}

var allMeasurementInfo = []map[Measurement]*unitInfo{volumeInfo, weightInfo}

var measurementLookup = make(map[string]*unitInfo)

var volumeOutput = make(map[System]byMeasurement)
var weightOutput = make(map[System]byMeasurement)

func extractOutputs(sys System, m map[Measurement]*unitInfo) byMeasurement {
	s := make([]*unitInfo, 0, len(m))
	for _, u := range m {
		if u.system == sys && u.outputType&SuppressOutput == 0 {
			s = append(s, u)
		}
	}
	result := byMeasurement(s)
	sort.Sort(sort.Reverse(result))
	return result
}

func init() {
	for _, m := range allMeasurementInfo {
		for measurement, info := range m {
			info.measurement = measurement

			for _, in := range info.in {
				measurementLookup[in] = info
			}
		}
	}

	volumeOutput[Metric] = extractOutputs(Metric, volumeInfo)
	volumeOutput[Imperial] = extractOutputs(Imperial, volumeInfo)

	weightOutput[Metric] = extractOutputs(Metric, weightInfo)
	weightOutput[Imperial] = extractOutputs(Imperial, weightInfo)
}

// 	}
// }

// need maps of accepted strings to the matching values

// Is there some kind of 'unit' interface that these both implement?
// Seems like we want a method that accepts a string and returns an
// insance of that unit
