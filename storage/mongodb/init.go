package mongodb

import (
	"btradoc/entities"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCollection struct {
	Name    string
	Indexes []mongo.IndexModel
}

var (
	ttlOnGoingTranslation int32 = 64800  // 18 hours
	ttlTemporaryTokens    int32 = 259200 // 3 days
	uniqueness            bool  = true
	noneLanguage                = "none"
)

var mongo_COLLECTIONS = [...]mongoCollection{
	{
		Name:    Datasets.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name:    Occitan.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: OnGoingTranslations.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "createdAt", Value: 1}},
				Options: &options.IndexOptions{ExpireAfterSeconds: &ttlOnGoingTranslation},
			},
		},
	},
	{
		Name:    SecretQuestions.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: TemporaryTokens.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "issuedAt", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, ExpireAfterSeconds: &ttlTemporaryTokens},
			},
		},
	},
	{
		Name:    Translations.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name:    TranslationsFiles.ColletionName(),
		Indexes: []mongo.IndexModel{},
	},
	{
		Name: Translations.ColletionName(),
		Indexes: []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
			{
				Keys:    bson.D{{Key: "refreshToken", Value: 1}},
				Options: &options.IndexOptions{Unique: &uniqueness, DefaultLanguage: &noneLanguage},
			},
		},
	},
}

// InitMongoDatabase initializes datas
func InitMongoDatabase() {
	occitan := openDialectsJSONFile("./data/occitan.json")
	secretQuestions := openSecretQuestions("./data/secret-questions.json")

	mongodb := NewMongoClient()
	createCollections(mongodb, mongo_COLLECTIONS[:])
	addOccitanDialects(mongodb, occitan)
	addSecretQuestions(mongodb, secretQuestions)
}

// createCollections creates all mongo collections
func createCollections(db *mongo.Database, collections []mongoCollection) {
	ctx := context.Background()

	for _, collection := range collections {
		err := db.CreateCollection(ctx, collection.Name)
		if err != nil {
			switch err.(type) {
			case mongo.CommandError:
				// continue if collection is already created
				continue
			default:
				log.Fatalln(err)
			}
			// if strings.Contains(err.Error(), "(NamespaceExists) Collection already exists") {
			// 	continue
			// }
		}

		coll := db.Collection(collection.Name, nil)
		coll.Indexes().CreateMany(ctx, collection.Indexes)
	}
}

// addOccitanDialects adds all occitan dialects to a specified Collection
func addOccitanDialects(db *mongo.Database, occitan []entities.OccitanJSONFile) {
	occitanColl := db.Collection("Occitan")
	count, err := occitanColl.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return
	}

	var totalDialect int
	for _, dialect := range occitan {
		totalDialect += len(dialect.Subdialects)
	}

	log.Printf("Total Occitan Dialects: %d (JSON) | %d (Documents)\n", totalDialect, count)

	ctx := context.Background()
	for _, dialect := range occitan {
		for _, subdialect := range dialect.Subdialects {
			result := occitanColl.FindOne(ctx, bson.D{{Key: "dialectName", Value: dialect.Dialect}, {Key: "subdialectName", Value: subdialect}})
			if result.Err() != nil {
				if result.Err() == mongo.ErrNoDocuments {
					// insert occitan dialect if it's not in documents
					if _, err := occitanColl.InsertOne(ctx, bson.D{
						{Key: "dialectName", Value: dialect.Dialect}, {Key: "subdialectName", Value: subdialect},
					}); err != nil {
						log.Println(err)
					}
				}

				continue
			}
		}
	}
}

// addSecretQuestions adds all secret questions to a specified Collection
func addSecretQuestions(db *mongo.Database, secretQuestions []string) {
	secretQuestionsColl := db.Collection("SecretQuestions")
	count, err := secretQuestionsColl.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Total secret questions: %d (JSON) | %d (Documents)\n", len(secretQuestions), count)

	ctx := context.Background()
	for _, secretQuestion := range secretQuestions {
		result := secretQuestionsColl.FindOne(ctx, bson.D{{Key: "question", Value: secretQuestion}})
		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				// insert secret question if it's not in documents
				if _, err := secretQuestionsColl.InsertOne(ctx, bson.D{{Key: "question", Value: secretQuestion}}); err != nil {
					log.Println(err)
				}
			}

			continue
		}
	}
}
