package translation

import (
	"btradoc/entities"
	"context"
	"errors"

	arangodb "github.com/arangodb/go-driver"
)

type Repository interface {
	GetDatasets(translatorID string, dialectName string, subdialectName string) (*[]entities.Dataset, error)
	InsertTranslations(translations []entities.Translation) error
}

type repository struct {
	ArangoDB arangodb.Database
}

func NewRepo(arangoDB arangodb.Database) Repository {
	return &repository{
		ArangoDB: arangoDB,
	}
}

func (r *repository) GetDatasets(translatorID string, dialectName string, subdialectName string) (*[]entities.Dataset, error) {
	query := `
		LET subdialectID = FIRST(
			FOR dialect IN dialects
				FILTER dialect.name == @dialectName
					FOR rds IN relDialectsSubdialects
						FILTER rds._from == dialect._id
						FOR subdialect IN subdialects
							FILTER subdialect._id == rds._to
							FILTER subdialect.name == @subdialectName
							return subdialect._id
		)
		LET exclude_dataset_already_translated = (
			FOR t IN translations
				FILTER t.subdialect == subdialectID
				return t.dataset
		)
		LET exclude_dataset_on_going_translations = (
			FOR ogt IN onGoingTranslations
				FILTER ogt.subdialect == subdialectID
				return ogt.dataset
		)
		FOR dataset IN datasets
			FILTER dataset._id NOT IN exclude_dataset_already_translated
			FILTER dataset._id NOT IN exclude_dataset_on_going_translations
			SORT RAND()
			LIMIT 5
			INSERT {
				subdialect: subdialectID,
				dataset: dataset._id,
				createdAt: DATE_TIMESTAMP(DATE_UTCTOLOCAL(DATE_NOW(), "Europe/Paris", true).local)
			} INTO onGoingTranslations
			return dataset
	`
	bindVars := map[string]interface{}{
		"dialectName":    dialectName,
		"subdialectName": subdialectName,
	}
	cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
	if err != nil {
		return nil, errors.New("A REMPLIR")
	}
	defer cursor.Close()

	datasets := []entities.Dataset{}
	for {
		var doc entities.Dataset
		_, err := cursor.ReadDocument(context.Background(), &doc)
		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, errors.New("")
		}
		datasets = append(datasets, doc)
	}

	return &datasets, nil
}

func (r *repository) InsertTranslations(translations []entities.Translation) error {
	for _, translation := range translations {
		// check if dataset is not already translated
		query := `
			FOR t IN translations
				FILTER t.dataset == @datasetID AND t.subdialect == @subdialectID
				return t
		`
		bindVars := map[string]interface{}{
			"datasetID":    translation.DatasetID,
			"subdialectID": translation.SubdialectID,
		}

		cursor, err := r.ArangoDB.Query(context.Background(), query, bindVars)
		if err != nil {
			return errors.New("A REMPLIR")
		}
		defer cursor.Close()

		if cursor.HasMore() {
			continue
		}

		// insert new translation
		query = `
			INSERT {
				oc: @oc, fr: @fr, en: @en, translator: @translatorID, dataset: @datasetID, subdialect: @subdialectID
			} INTO translations
		`
		bindVars = map[string]interface{}{
			"oc":           translation.Oc,
			"fr":           translation.Fr,
			"en":           translation.En,
			"translatorID": translation.TranslatorID,
			"datasetID":    translation.DatasetID,
			"subdialectID": translation.SubdialectID,
		}

		_, err = r.ArangoDB.Query(context.Background(), query, bindVars)
		if err != nil {
			return errors.New("A REMPLIR")
		}

		// Remove onGoingTranslations
		query = `
			FOR ogt IN onGoingTranslations
				FILTER ogt.subdialect == @subdialectID
				FILTER ogt.dataset == @datasetID
				REMOVE ogt IN onGoingTranslations
				RETURN OLD
		`
		bindVars = map[string]interface{}{
			"subdialectID": translation.SubdialectID,
			"datasetID":    translation.DatasetID,
		}

		_, err = r.ArangoDB.Query(context.Background(), query, bindVars)
		if err != nil {
			return errors.New("A REMPLIR")
		}
	}
	return nil
}
