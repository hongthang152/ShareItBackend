package repository

import (
	"context"

	"github.com/hongthang152/ShareItBackend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepo interface {
	FetchAll(ctx context.Context) ([]models.File, error)
	FetchAllOneWeekExpired(ctx context.Context) ([]models.File, error)
	DeleteAllOneWeekExpired(ctx context.Context) (*mongo.DeleteResult, error)
	FetchByPIN(ctx context.Context, pin string) (*models.File, error)
	Create(ctx context.Context, file models.File) (*mongo.InsertOneResult, error)
}
