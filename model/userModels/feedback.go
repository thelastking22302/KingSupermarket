package usermodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	Id           primitive.ObjectID `bson:"_id"`
	Feedback_id  *string            `json:"feedback_id" bson:"feedback_id"`
	Name         string             `json:"name" bson:"name" validate:"required"`
	Phone_Number string             `json:"phone_number" bson:"phone" validate:"required"`
	Note         string             `json:"note" bson:"note" validate:"required"`
	Created_at   *time.Time         `json:"created_at" bson:"created_at"`
	Updated_at   *time.Time         `json:"updated_at" bson:"updated_at"`
}
