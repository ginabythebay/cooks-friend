package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Section struct {
	Name        string
	Ingredients []Ingredient
	Steps       []string
}

type Ingredient struct {
	Item   string
	fields []string
}

func (i *Ingredient) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&i.fields)
	if err == nil {
		// TODO(gina) error check that the item exists, etc
		i.Item = i.fields[0]
	}
	return err
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
			for idx, fld := range i.fields {
				if idx != 0 {
					log.Print(fld)
					Parse(fld)
				}
			}
		}
	}
}
