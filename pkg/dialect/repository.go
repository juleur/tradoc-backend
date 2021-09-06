package dialect

import (
	"btradoc/entities"
	"context"
	"errors"

	arangodb "github.com/arangodb/go-driver"
)

type Repository interface {
	GetDialectPermissions(translatorID string) (*[]entities.DialectSubdialectDocument, error)
	GetDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error)
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangoDB arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangoDB,
	}
}

func (r *repository) GetDialectPermissions(translatorID string) (*[]entities.DialectSubdialectDocument, error) {
	query := `
		FOR dp IN dialectPermissions
			FILTER dp.translator == @translatorID
			LET subdialect = FIRST(
				FOR s IN subdialects
					FILTER s._id == dp.subdialect
					return s.name
			)
			LET dialect = FIRST(
				FOR rds IN relDialectsSubdialects
					FILTER rds._to == dp.subdialect
						FOR d IN dialects
							FILTER d._id == rds._from
							return d.name
			)
			return {
				dialect: dialect,
				subdialect: subdialect
			}
	`
	bindVars := map[string]interface{}{
		"translatorID": translatorID,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New("A REMPLIR")
	}
	defer cursor.Close()

	dialectSubdialectDoc := []entities.DialectSubdialectDocument{}
	for {
		var doc entities.DialectSubdialectDocument
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errors.New("")
		}
		dialectSubdialectDoc = append(dialectSubdialectDoc, doc)
	}

	return &dialectSubdialectDoc, nil
}

func (r *repository) GetDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error) {
	query := `
		FOR d IN dialects
			LET subdialects = (
				FOR rds IN relDialectsSubdialects
				FILTER rds._from == d._id
				FOR s IN subdialects
					FILTER s._id == rds._to
					LET totalTransl = LENGTH(FOR t IN translations FILTER t.subdialect == s._id RETURN 1)
					LET totalTranslByTranslator = LENGTH(
						FOR t IN translations
							FILTER t.subdialect == s._id AND t.translator == @translatorID
							RETURN 1
					)
				return {
					subdialect_id: s._id,
					subdialect: s.name,
					total_translations: totalTransl,
					total_translations_by_translator: totalTranslByTranslator
				}
			)
			return {
				dialect_id: d._id,
				dialect: d.name,
				subdialects: subdialects
			}
	`
	bindVars := map[string]interface{}{
		"translatorID": translatorID,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New("A REMPLIR")
	}
	defer cursor.Close()

	dialectSubdialects := []entities.DialectSubdialects{}
	for {
		var doc entities.DialectSubdialects
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errors.New("")
		}
		dialectSubdialects = append(dialectSubdialects, doc)
	}

	return &dialectSubdialects, nil
}
