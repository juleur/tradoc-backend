package dialect

import (
	"btradoc/entities"
	"btradoc/pkg"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error)
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

func (r *repository) GetDialectsSubdialect(translatorID string) (*[]entities.DialectSubdialects, error) {
	occitanColl := r.MongoDB.Collection("Occitan")
	ID, err := primitive.ObjectIDFromHex(translatorID)
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Translations"},
			{Key: "pipeline", Value: []interface{}{bson.D{{Key: "$match", Value: bson.D{{Key: "fr", Value: bson.D{{Key: "$exists", Value: true}}}}}}}},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "occitan"},
			{Key: "as", Value: "totalTranslated"},
		}},
	}
	project1Stage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "dialectName", Value: 1},
			{Key: "subdialectName", Value: 1},
			{Key: "totalTranslated", Value: 1},
			{Key: "totalTranslatedByTranslator", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$totalTranslated"},
					{Key: "as", Value: "tt"},
					{Key: "cond", Value: bson.D{{Key: "$eq", Value: []interface{}{"$$tt.translator", ID}}}},
				},
				}}},
		}},
	}
	project2Stage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "dialectName", Value: 1},
			{Key: "subdialectName", Value: 1},
			{Key: "totalTranslated", Value: bson.D{{Key: "$size", Value: "$totalTranslated"}}},
			{Key: "totalTranslatedByTranslator", Value: bson.D{{Key: "$size", Value: "$totalTranslatedByTranslator"}}},
		}},
	}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$dialectName"},
			{Key: "subdialects", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "name", Value: "$$ROOT.subdialectName"},
				{Key: "totalTranslated", Value: "$totalTranslated"},
				{Key: "totalTranslatedByTranslator", Value: "$totalTranslatedByTranslator"},
			}}}},
		}},
	}
	project3Stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "dialect", Value: "$_id"}, {Key: "subdialects", Value: 1}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "dialect", Value: 1}}}}

	ctx := context.Background()
	cursor, err := occitanColl.Aggregate(ctx, mongo.Pipeline{lookupStage, project1Stage, project2Stage, groupStage, project3Stage, sortStage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var result []entities.DialectSubdialects
	if err = cursor.All(ctx, &result); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	} else if len(result) == 0 {
		return nil, &pkg.DBError{
			Code:    404,
			Message: pkg.ErrDialectNotFound,
			Wrapped: err,
		}
	}

	return &result, nil
}
