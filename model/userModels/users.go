package usermodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id           primitive.ObjectID `bson:"_id"`
	User_id      string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name         *string            `json:"name" validate:"required" bson:"name"`
	Avatar       *string            `json:"avatar" validate:"required" bson:"avatar"`
	Age          int                `json:"age" validate:"required" bson:"age"`
	Email        *string            `json:"email" validate:"required" bson:"email"`
	Password     string             `json:"password" validate:"required" bson:"password"`
	Role         string             `json:"role"  bson:"role"`
	Phone_Number *string            `json:"phone_number" validate:"required" bson:"phone_number"`
	Created_At   *time.Time         `json:"created_at"  bson:"created_at"`
	Updated_At   *time.Time         `json:"updated_at" bson:"updated_at"`
}
