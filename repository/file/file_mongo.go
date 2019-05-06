package file

import (
	"context"
	"github.com/hongthang152/ShareItBackend/models"
	"github.com/hongthang152/ShareItBackend/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func NewMongoRepository(Client *mongo.Client) repository.FileRepo {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	return &repo{
		Client:     Client,
		Collection: Client.Database("shareit").Collection("file"),
	}
}

type repo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func (mongoRepo *repo) FetchByPIN(ctx context.Context, pin string) (*models.File, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	filter := bson.M{"pin": pin}
	var result models.File
	err := mongoRepo.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (mongoRepo *repo) Create(ctx context.Context, file models.File) (*mongo.InsertOneResult, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	return mongoRepo.Collection.InsertOne(ctx, bson.M{
		"name":      file.Name,
		"pin":       file.Pin,
		"url":       file.Url,
		"createdAt": time.Now().UTC().Unix(),
	})
}

func (mongoRepo *repo) FetchAll(ctx context.Context) ([]models.File, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	cur, err := mongoRepo.Collection.Find(ctx, bson.M{})
	var results []models.File

	for cur.Next(ctx) {
		var result models.File
		cur.Decode(&result)
		results = append(results, result)
	}

	return results, err
}

func (mongoRepo *repo) FetchAllOneWeekExpired(ctx context.Context) ([]models.File, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	cur, err := mongoRepo.Collection.Find(ctx, bson.M{
		"createdAt": bson.M{
			"$lt": time.Now().UTC().Unix() - 604800,
		},
	})
	defer cur.Close(ctx)
	var results []models.File

	for cur.Next(ctx) {
		var result models.File
		cur.Decode(&result)
		results = append(results, result)
	}

	return results, err
}

func (mongoRepo *repo) DeleteAllOneWeekExpired(ctx context.Context) (*mongo.DeleteResult, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	return mongoRepo.Collection.DeleteMany(ctx, bson.M{
		"createdAt": bson.M{
			"$lt": time.Now().UTC().Unix() - 604800,
		},
	})
}
