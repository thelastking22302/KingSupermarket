package marketmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItems struct {
	Id            primitive.ObjectID `bson:"_id"`
	Order_Item_Id *string            `json:"order_item_id,omitempty" bson:"order_item_id,omitempty"`
	Order_Id      *string            `json:"order_id" bson:"order_id" validate:"required"`
	Product_id    *string            `json:"product_id" bson:"product_id" validate:"required"`
	Quantity      int                `json:"quantity" bson:"quantity" validate:"required"`
	Price         float64            `json:"price" bson:"price" validate:"required"`
	Created_at    *time.Time         `json:"created_at" bson:"created_at"`
	Updated_at    *time.Time         `json:"updated_at" bson:"updated_at"`
}
