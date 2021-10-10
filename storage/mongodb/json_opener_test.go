package mongodb

import "testing"

func TestOpenDialectsJSONFile(t *testing.T) {
	occitan := openDialectsJSONFile("./../../data/occitan.json")
	t.Log(occitan)
}

func TestOpenSecretQuestions(t *testing.T) {
	secretQuestions := openSecretQuestions("./../../data/secret-questions.json")
	t.Log(secretQuestions)
}
