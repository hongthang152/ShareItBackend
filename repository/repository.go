package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

type FileRepo interface {
	FetchByPIN(ctx context.Context, pin int64) (*models.File, error)
}