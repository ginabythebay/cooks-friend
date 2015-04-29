package main

import "fmt"

type Recipe struct {
	Title    string
	Sections []Section
}

func (r *Recipe) ShoppingList() ([]Ingredient, error) {
	capacity := 0
	for _, s := range r.Sections {
		capacity += len(s.Ingredients)
	}

	// this will have more capacity than needed if the same ingredient appears twice
	result := make([]Ingredient, 0, capacity)
	have := make(map[string]int) // from name to location
	for _, s := range r.Sections {
		for _, i := range s.Ingredients {
			offset, found := have[i.Item]
			if found {
				if err := result[offset].Merge(&i); err != nil {
					return nil, err
				}
			} else {
				result = append(result, i)
				have[i.Item] = len(result) - 1
			}
		}
	}
	return result, nil
}

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

func (i *Ingredient) Merge(other *Ingredient) error {
	if i.Item != other.Item {
		return fmt.Errorf("Cannot add ingredient %q to %q)", i.Item, other.Item)
	}
	if len(i.Measurements) != len(other.Measurements) {
		return fmt.Errorf("Cannot merge ingredients.  %q has %d measurements while %q has %d measurements.",
			i.Item, len(i.Measurements), other.Item, len(other.Measurements))
	}
	// we add up the new measurements first before assigning them so
	// we don't modify i if there are any errors
	temp := make([]Measurement, 0, len(i.Measurements))
	for index, m := range i.Measurements {
		m, err := m.Add(other.Measurements[index])
		if err != nil {
			return err
		}
		temp = append(temp, m)
	}
	for index, m := range temp {
		i.Measurements[index] = m
	}
	return nil
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
