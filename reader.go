package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type Section struct {
	Name        string
	Ingredients []Ingredient
	Steps       []string
}

// We want to capture both the weight and the volume if both are present.
// If weight is present, we will prefer to display that in the recipe.
//
// When displaying a shopping list, we can combine ingredients that
// have both the same name nad the same types of measurements
// (e.g. both volume or both weight or both weight and volume).
type Ingredient struct {
	Item string
	// what do we want here?  Map from type to Measurement?  called-out volume and weight?
	Measurements []Measurement
}

func (i *Ingredient) UnmarshalYAML(unmarshal func(interface{}) error) error {
	fields := make([]string, 0)
	if err := unmarshal(&fields); err != nil {
		return err
	}
	// TODO(gina) verify we have an item and at least one measurement
	i.Item = fields[0]
	i.Measurements = make([]Measurement, len(fields)-1)
	for idx, f := range fields {
		if idx == 0 {
			continue
		}
		var err error
		if i.Measurements[idx-1], err = Parse(f); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	b, err := ioutil.ReadFile("/home/gina/go/src/github.com/ginabythebay/cooks-friend/recipes/whole-wheat-rustic-italian-bread.yml")
	if err != nil {
		log.Fatal(err)
	}
	var recipe []Section
	err = yaml.Unmarshal(b, &recipe)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Read recipe: %+v", recipe)

	for _, s := range recipe {
		for _, i := range s.Ingredients {
			tokens := make([]string, len(i.Measurements))
			for idx, m := range i.Measurements {
				tokens[idx] = m.Output(Imperial)
			}
			log.Println(i.Item, ": ", strings.Join(tokens, " or "))
		}
	}
}
