package mongodb

import "testing"

func TestCreateCollections(t *testing.T) {
	mongodb := NewMongoClient()
	createCollections(mongodb, mongo_COLLECTIONS[:])
}

func TestAddOccitanDialects(t *testing.T) {
	occitan := openDialectsJSONFile("./../../data/occitan.json")
	mongodb := NewMongoClient()
	addOccitanDialects(mongodb, occitan)
}

func TestAddSecretQuestions(t *testing.T) {
	secretQuestions := openSecretQuestions("./../../data/secret-questions.json")
	mongodb := NewMongoClient()
	addSecretQuestions(mongodb, secretQuestions)
}
