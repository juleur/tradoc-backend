package entities

import (
	arangodb "github.com/arangodb/go-driver"
)

type Collection struct {
	Name string
	Type arangodb.CollectionType
}

type Colletions []Collection

var COLLECTIONS = []Collection{
	{
		Name: "datasets",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "dialectPermissions",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "dialects",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "onGoingTranslations",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "refreshTokens",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "relDialectsSubdialects",
		Type: arangodb.CollectionTypeEdge,
	},
	{
		Name: "subdialects",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "translations",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "translationsFiles",
		Type: arangodb.CollectionTypeDocument,
	},
	{
		Name: "translators",
		Type: arangodb.CollectionTypeDocument,
	},
}
