package translation

import (
	"btradoc/entities"
	"btradoc/pkg"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetDatasets(fullDialect string) (*[]entities.Dataset, error)
	AddNewFullDialectToDataset(translations []*entities.Translation) error
	InsertTranslations(translations []*entities.Translation) error
	GetTotalOnGoingTranslation(fullDialect, translatorID string) (int, error)
	InsertDatasetsOnGoingTranslations(fullDialect, translatorID string, datasets *[]entities.Dataset) error
	RemoveDatasetsOnGoingTranslations(translations []*entities.Translation) error
	GetTranslationsFiles() (*[]entities.TranslationFile, error)
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

// GetDatasets gets datasets according the dialect and subdialect wanted
func (r *repository) GetDatasets(fullDialect string) (*[]entities.Dataset, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "translatedIn", Value: bson.D{{Key: "$nin", Value: []interface{}{fullDialect}}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "OnGoingTranslations"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "dataset"}, {Key: "as", Value: "join"}}}}
	match2Stage := bson.D{{Key: "$match", Value: bson.D{{Key: "join", Value: bson.D{{Key: "$size", Value: 0}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "sentence", Value: 1}}}}
	sampleStage := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}}

	datasetsColl := r.MongoDB.Collection("Datasets")
	ctx := context.Background()
	cursor, err := datasetsColl.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, match2Stage, projectStage, sampleStage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var datasets []entities.Dataset
	if err = cursor.All(ctx, &datasets); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return &datasets, nil
}

func (r *repository) AddNewFullDialectToDataset(translations []*entities.Translation) error {
	datasetsColl := r.MongoDB.Collection("Datasets")

	for _, translation := range translations {
		if _, err := datasetsColl.UpdateOne(context.Background(), bson.M{"dataset": translation.DatasetID}, bson.D{{Key: "$push", Value: bson.D{{Key: "translatedIn", Value: translation.FullDialect}}}}); err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}
	}

	return nil
}

// InsertTranslations inserts datasets translated in Occitan dialect -> French and/or Occitan dialect -> English
func (r *repository) InsertTranslations(translations []*entities.Translation) error {
	translationColl := r.MongoDB.Collection("Translations")

	var translationDocs []interface{}
	for _, translation := range translations {
		translatorID, err := primitive.ObjectIDFromHex(translation.TranslatorID)
		if err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}
		datasetID, err := primitive.ObjectIDFromHex(translation.DatasetID)
		if err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		// english translation
		if len(translation.En) == 0 {
			doc := bson.D{
				{Key: "oc", Value: translation.Oc},
				{Key: "en", Value: translation.En},
				{Key: "translator", Value: translatorID},
				{Key: "dataset", Value: datasetID},
				{Key: "occitan", Value: translation.FullDialect},
			}
			translationDocs = append(translationDocs, doc)
		}
		// french translation
		doc := bson.D{
			{Key: "oc", Value: translation.Oc},
			{Key: "fr", Value: translation.Fr},
			{Key: "translator", Value: translatorID},
			{Key: "dataset", Value: datasetID},
			{Key: "occitan", Value: translation.FullDialect},
		}
		translationDocs = append(translationDocs, doc)
	}

	if _, err := translationColl.InsertMany(context.Background(), translationDocs); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}

// GetTotalOnGoingTranslation returns how many translations are being translated by given translator in order to prevent spamming
func (r *repository) GetTotalOnGoingTranslation(fullDialect, translatorID string) (int, error) {
	translatorObjectID, err := primitive.ObjectIDFromHex(translatorID)
	if err != nil {
		return 0, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ctx := context.Background()

	onGoingTranslations := r.MongoDB.Collection("OnGoingTranslations")

	cursor, err := onGoingTranslations.Find(ctx, primitive.D{{Key: "translator", Value: translatorObjectID}, {Key: "occitan", Value: fullDialect}})
	if err != nil {
		return 0, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var ogtDocuments []bson.M
	if err = cursor.All(ctx, &ogtDocuments); err != nil {
		return 0, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return len(ogtDocuments), nil
}

// InsertDatasetsOnGoingTranslations inserts datasets that's being translated in order to prevent duplications
func (r *repository) InsertDatasetsOnGoingTranslations(fullDialect, translatorID string, datasets *[]entities.Dataset) error {
	translatorObjectID, err := primitive.ObjectIDFromHex(translatorID)
	if err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ogtColl := r.MongoDB.Collection("OnGoingTranslations")

	var datasetDocs []interface{}
	for _, d := range *datasets {
		ID, err := primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}
		doc := bson.D{
			{Key: "fullDialect", Value: "fmt.Sprintf"},
			{Key: "occitan", Value: fullDialect},
			{Key: "dataset", Value: ID},
			{Key: "translator", Value: translatorObjectID},
			{Key: "createdAt", Value: time.Now()},
		}

		datasetDocs = append(datasetDocs, doc)
	}

	if _, err = ogtColl.InsertMany(context.Background(), datasetDocs); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}

// RemoveDatasetsOnGoingTranslations deletes dataset that has been translated
func (r *repository) RemoveDatasetsOnGoingTranslations(translations []*entities.Translation) error {
	onGoingTranslationsColl := r.MongoDB.Collection("OnGoingTranslations")

	for _, translation := range translations {
		datasetIDObject, err := primitive.ObjectIDFromHex(translation.DatasetID)
		if err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		if _, err := onGoingTranslationsColl.DeleteOne(context.Background(), bson.D{{Key: "dataset", Value: datasetIDObject}, {Key: "occitan", Value: translation.FullDialect}}); err != nil {
			return &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}
	}

	return nil
}

// GetTranslationsFiles fetches translation filenames (occitan->french & occitan->english)
func (r *repository) GetTranslationsFiles() (*[]entities.TranslationFile, error) {
	translationsFilesColl := r.MongoDB.Collection("TranslationsFiles")

	ctx := context.Background()

	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$dialectName"}, {Key: "subdialects", Value: bson.D{{Key: "$push", Value: bson.D{{Key: "name", Value: "$$ROOT.subdialectName"}, {Key: "files", Value: "$$ROOT.files"}}}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "dialect", Value: "$_id"}, {Key: "subdialects", Value: "$subdialects"}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "dialect", Value: 1}}}}

	cursor, err := translationsFilesColl.Aggregate(ctx, mongo.Pipeline{groupStage, projectStage, sortStage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var tfs []entities.TranslationFile
	if err = cursor.All(ctx, &tfs); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return &tfs, nil
}
