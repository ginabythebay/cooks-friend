package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	b, err := ioutil.ReadFile("/home/gina/go/src/github.com/ginabythebay/cooks-friend/testdata/whole-wheat-rustic-italian-bread.yml")
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
