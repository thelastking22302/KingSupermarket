package marketmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id"`
	Product_id  *string            `json:"product_id" bson:"product_id"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Image       *string            `json:"image" bson:"image" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	Price       float64            `json:"price" bson:"price" validate:"required"`
	Stock       int                `json:"stock" bson:"stock" validate:"required"`
	Status      string             `json:"status" bson:"status" validate:"required"`
	Category_id *string            `json:"category_id" bson:"category_id" validate:"required"`
	Created_at  *time.Time         `json:"created_at" bson:"created_at"`
	Updated_at  *time.Time         `json:"updated_at" bson:"updated_at"`
}
