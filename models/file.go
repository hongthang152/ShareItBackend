package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" binding:"required"`
	Name      string             `json:"name" bson:"name,omitempty" binding:"required"`
	Pin       string             `json:"pin" bson:"pin,omitempty" binding:"required"`
	Url       string             `json:"url" bson:"url,omitempty" binding:"required"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt,omitempty" binding:"required"`
}
