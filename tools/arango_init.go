package tools

import (
	"btradoc/entities"
	"context"
	"fmt"
	"log"

	arangodb "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func CreateCollection(arangoDB arangodb.Database, collections entities.Colletions) error {
	for _, col := range collections {
		if ok, err := arangoDB.CollectionExists(context.Background(), col.Name); ok || err != nil {
			if err != nil {
				return err
			}
			return fmt.Errorf("“%s“ collection already added", col.Name)
		}

		options := &arangodb.CreateCollectionOptions{
			Type: col.Type,
		}
		if _, err := arangoDB.CreateCollection(context.Background(), col.Name, options); err != nil {
			return err
		}
	}

	return nil
}

func AddDialectsDocuments(arangoDB arangodb.Database, dialects []entities.Dialect) {
	dialectCol, err := arangoDB.Collection(context.Background(), "dialects")
	if err != nil {
		log.Fatalln(err)
	}

	subdialectCol, err := arangoDB.Collection(context.Background(), "subdialects")
	if err != nil {
		log.Fatalln(err)
	}

	relDialectsSubdialectsCol, err := arangoDB.Collection(context.Background(), "relDialectsSubdialects")
	if err != nil {
		log.Fatalln(err)
	}

	for _, dialect := range dialects {
		document := struct {
			Name string `json:"name"`
		}{
			Name: dialect.Name,
		}
		metaDialect, err := dialectCol.CreateDocument(context.Background(), document)
		if err != nil {
			log.Fatalln(err)
		}

		for _, subdialect := range dialect.Subdialects {
			document = struct {
				Name string `json:"name"`
			}{
				Name: subdialect,
			}
			metaSubdialect, err := subdialectCol.CreateDocument(context.Background(), document)
			if err != nil {
				log.Fatalln(err)
			}

			edge := arangodb.EdgeDocument{
				From: metaDialect.ID,
				To:   metaSubdialect.ID,
			}
			if _, err = relDialectsSubdialectsCol.CreateDocument(context.Background(), edge); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func ArangoDBConnection() arangodb.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529"},
	})
	if err != nil {
		log.Fatalln(err)
	}

	c, err := arangodb.NewClient(arangodb.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatalln(err)
	}

	db, err := c.Database(context.Background(), "oc")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
