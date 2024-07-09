package marketmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	Id          primitive.ObjectID `bson:"_id"`
	Category_id *string            `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Created_at  *time.Time         `json:"created_at" bson:"created_at"`
	Updated_at  *time.Time         `json:"updated_at" bson:"updated_at"`
}
