package file

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func NewMongoRepository(Client *mongo.Client) Repository {
	return &repo{
		Client: Client,
	}
}

type repo struct {
	Client *mongo.Client
}

func (mongoRepo *repo) FetchByPIN(ctx context.Context, pin string, args ...interface{}) (*models.File, error) {
	collection := mongoRepo.Client.Database("ShareIt").Collection("file")
	filter := bson.M{"pin": pin}
	var result models.File
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return &result, err
}
