package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_URI      string = "mongodb://localhost:27017"
	MONGO_USERNAME string = "occitan" //  /!\ not safe
	MONGO_PASSWORD string = "naticco" //  /!\ not safe
	MONGO_DB_NAME  string = "oc"
)

func NewMongoClient() *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MONGO_URI).SetDirect(true).SetAuth(options.Credential{
		Username: MONGO_USERNAME,
		Password: MONGO_PASSWORD,
	}))

	db := client.Database(MONGO_DB_NAME)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
