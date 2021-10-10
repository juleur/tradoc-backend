package dialect

import (
	"btradoc/storage/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchDialectsSubdialect(t *testing.T) {
	db := mongodb.NewMongoClient()
	translatorRepo := NewRepo(db)
	translatorService := NewService(translatorRepo)

	translatorID := "6148c3f1ba78b40cdeb49289"
	result, err := translatorService.FetchDialectsSubdialect(translatorID)
	assert.Nil(t, err)
	t.Log(result)
}
