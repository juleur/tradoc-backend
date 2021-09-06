package translation

import (
	"btradoc/entities"
	"btradoc/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchSentencesToTranslate(t *testing.T) {
	db := tools.ArangoDBConnection()

	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translatorID := "translators/59819"
	dialectName := "provençau"
	subdialectName := "brivadés"
	datasets, err := translationService.FetchSentencesToTranslate(translatorID, dialectName, subdialectName)
	assert.Nil(t, err)
	t.Log(datasets)
}

func TestAddTranslations(t *testing.T) {
	db := tools.ArangoDBConnection()

	translationRepo := NewRepo(db)
	translationService := NewService(translationRepo)

	translations := []entities.Translation{
		{
			Oc:           "Cossí l'aimatz lo cafè ? Sucrat ?",
			Fr:           "Comment l'aimez-vous le café ? Sucré ?",
			En:           "How do you like your coffee ? Sugary ?",
			TranslatorID: "translators/104849",
			DatasetID:    "datasets/110420",
			SubdialectID: "subdialects/15005",
		},
		{
			Oc:           "Volètz dormir una estona per vos pausar ?",
			Fr:           "Voulez-vous dormir un moment pour vous reposer ?",
			En:           "Do you want to sleep for a while to rest ?",
			TranslatorID: "translators/82678",
			DatasetID:    "datasets/110422",
			SubdialectID: "subdialects/14997",
		},
	}

	err := translationService.AddTranslations(translations)
	assert.Nil(t, err)
}
