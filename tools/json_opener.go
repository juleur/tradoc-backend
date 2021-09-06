package tools

import (
	"btradoc/entities"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func OpenDialectsJSONFile() []entities.Dialect {
	jsonFile, err := os.Open("./data/dialects-subdialects.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	dialects := []entities.Dialect{}
	err = json.Unmarshal(byteValue, &dialects)
	if err != nil {
		log.Fatalln(err)
	}

	return dialects
}
