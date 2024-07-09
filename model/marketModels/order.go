package marketmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id           primitive.ObjectID `bson:"_id"`
	Order_Id     *string            `json:"order_id,omitempty" bson:"order_id,omitempty"`
	User_Id      *string            `json:"user_id" bson:"user_id"`
	Address      string             `json:"address" bson:"address" validate:"required"`
	Phone_Number *string            `json:"phone_number" bson:"phone_number" validare:"required"`
	Total_amount float64            `json:"total_amount" bson:"total_amount"`
	Status       string             `json:"status" bson:"status" validate:"eq=processing|eq=success|eq=defeat"`
	Notes        *string            `json:"notes" bson:"notes"`
	Order_day    *time.Time         `json:"order_day" bson:"order_day"`
	Created_at   *time.Time         `json:"created_at" bson:"created_at"`
	Updated_at   *time.Time         `json:"updated_at" bson:"updated_at"`
}
