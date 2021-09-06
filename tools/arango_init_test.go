package tools

import (
	"btradoc/entities"
	"testing"
)

func TestCollectionsInit(t *testing.T) {
	db := ArangoDBConnection()

	dialects := OpenDialectsJSONFile()

	CreateCollection(db, entities.COLLECTIONS)
	AddDialectsDocuments(db, dialects)
}
