package dialect

import (
	"btradoc/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDialectPermissions(t *testing.T) {
	db := tools.ArangoDBConnection()

	dialectRepo := NewRepo(db)
	dialectService := NewService(dialectRepo)

	translatorID := "translators/33460"
	result, err := dialectService.FindDialectPermissions(translatorID)
	assert.Nil(t, err)
	t.Log(result)

	translatorID = "translators/64926"
	result, err = dialectService.FindDialectPermissions(translatorID)
	assert.Nil(t, err)
	t.Log(result)

	translatorID = "translators/64922"
	result, err = dialectService.FindDialectPermissions(translatorID)
	assert.Nil(t, err)
	t.Log(result)
}

func TestFetchDialectsSubdialect(t *testing.T) {
	db := tools.ArangoDBConnection()

	dialectRepo := NewRepo(db)
	dialectService := NewService(dialectRepo)

	translatorID := "translators/33460"
	result, err := dialectService.FetchDialectsSubdialect(translatorID)
	assert.Nil(t, err)
	t.Log(result)
}
