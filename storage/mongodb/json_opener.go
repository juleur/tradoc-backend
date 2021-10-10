package mongodb

import (
	"btradoc/entities"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// openDialectsJSONFile opens occitan dialects json file where all subdialects are ordered by dialect
func openDialectsJSONFile(path string) []entities.Occitan {
	byteValue := opener(path)

	dialects := []entities.Occitan{}
	err := json.Unmarshal(byteValue, &dialects)
	if err != nil {
		log.Fatalln(err)
	}

	return dialects
}

// openSecretQuestions opens a json file where secret questions are saved in order to use them for reset password
func openSecretQuestions(path string) []string {
	byteValue := opener(path)

	var secretQuestions []string
	err := json.Unmarshal(byteValue, &secretQuestions)
	if err != nil {
		log.Fatalln(err)
	}

	return secretQuestions
}

func opener(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	return byteValue
}
