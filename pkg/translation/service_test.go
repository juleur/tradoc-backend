package translation

import (
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchSentencesToTranslate(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	fullDialect := "auvernhat_estandard"
	datasets, err := translationService.FetchDatasets(fullDialect)
	assert.Nil(t, err)
	t.Logf("%+v\n", datasets)
}

// func TestAddTranslations(t *testing.T) {
// 	db := tools.ArangoDBConnection()

// 	translationRepo := NewRepo(db)
// 	translationService := NewService(translationRepo)

// 	translations := []*entities.Translation{
// 		{
// 			Oc: "Es de Canadà mai naissèt en Anglatèrra",
// 			Fr: "Je suis au Canada mais je suis né en Angleterre",
// 			// En:           "I'm in Canada but I was born in England",
// 			TranslatorID: "translators/3370",
// 			DatasetID:    "datasets/9969",
// 			FullDialect:  "auvernhat-estandard",
// 		},
// 		{
// 			Oc:           "E tanben parles fòrça ben",
// 			Fr:           "Et aussi tu parles assez bien",
// 			En:           "And also you talk pretty good",
// 			TranslatorID: "translators/3370",
// 			DatasetID:    "datasets/9975",
// 			FullDialect:  "auvernhat-estandard",
// 		},
// 	}

// 	err := translationService.AddTranslations(translations)
// 	assert.Nil(t, err)
// }

func TestGetTotalOnGoingTranslation(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translatorID := "6148c3f1ba78b40cdeb49288"
	fullDialect := "auvernhat_estandard"

	counter, err := translationService.FetchTotalOnGoingTranslations(fullDialect, translatorID)
	assert.Nil(t, err)
	t.Logf("total: %d", counter)
}

func TestInsertDatasetsOnGoingTranslations(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	fullDialect := "auvernhat_estandard"

	datasets, err := translationService.FetchDatasets(fullDialect)
	assert.Nil(t, err)

	translatorID := "6148c3f1ba78b40cdeb49288"

	err = translationService.AddOnGoingTranslations(fullDialect, translatorID, datasets)
	assert.Nil(t, err)
}

func TestGetTranslationsFiles(t *testing.T) {
	db := mongodb.NewMongoClient()
	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translationsFiles, err := translationService.FetchTranslationsFiles()
	assert.Nil(t, err)
	t.Logf("%+v\n", *translationsFiles)
}
