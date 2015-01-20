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

// TODO(gina) look at making this a struct with custom marshall/unmarshall handling
type Ingredient []string

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
}
