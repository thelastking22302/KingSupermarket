package repomarketiml

import (
	"context"
	"errors"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDbOrderItems(db *mongo.Client) *Db {
	return &Db{db: db}
}

func (db *Db) collectionOrderItems() *mongo.Collection {
	return db.db.Database("KingSupermarket").Collection("orderItems")
}

func (db *Db) CreateOrderItems(ctx context.Context, data *marketmodels.OrderItems) error {
	filter := bson.M{"order_id": data.Order_Id}
	cnt, err := db.collectionOrder().CountDocuments(ctx, filter)
	if err != nil || cnt == 0 {
		return errors.New("orderItems not found")
	}
	var dataOrder marketmodels.Order
	if err := db.collectionOrder().FindOne(ctx, filter).Decode(&dataOrder); err != nil {
		return errors.New("orderItems cannot be decoded")
	}
	filterProduct := bson.M{"product_id": data.Product_id}
	var productOrder marketmodels.Product
	if err := db.dbProducts().FindOne(ctx, filterProduct).Decode(&productOrder); err != nil {
		return errors.New("product can't on OrderItems")
	}
	result, err := db.collectionOrderItems().InsertOne(ctx, &data)
	if err != nil {
		return errors.New("orderItems faild")
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("orderItems invalid")
	}
	data.Id = oid
	return nil
}

func (db *Db) GetOrderItems(ctx context.Context, id string) (*marketmodels.OrderItems, error) {
	filter := bson.M{"order_item_id": id}
	var dataOrderItems marketmodels.OrderItems
	if err := db.collectionOrderItems().FindOne(ctx, filter).Decode(&dataOrderItems); err != nil {
		return nil, errors.New("orderItems invalid")
	}
	return &dataOrderItems, nil
}
func (db *Db) DeleteOrderItems(ctx context.Context, id string) error {
	filter := bson.M{"order_item_id": id}
	dataDel, err := db.collectionOrderItems().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if dataDel.DeletedCount == 0 {
		return errors.New("orderItems not found")
	}
	return nil
}
func (db *Db) UpdateOrderItems(ctx context.Context, id string, data *marketmodels.OrderItems) error {
	upd := bson.D{
		{
			Key: "$set", Value: bson.M{
				"product_id": data.Product_id,
				"quantity":   data.Quantity,
				"price":      data.Price,
				"updated_at": data.Updated_at,
			},
		},
	}
	filter := bson.M{"order_item_id": id}
	cnt, err := db.collectionOrderItems().UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	if cnt.ModifiedCount == 0 {
		return err
	}
	return nil
}
func (db *Db) GetListOrderItems(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.OrderItems, error) {
	pagging.Process()
	opts := options.Find().SetSkip((int64((pagging.Page - 1) * pagging.Limit))).SetLimit(int64(pagging.Limit))
	cur, err := db.collectionOrderItems().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var listItems []marketmodels.OrderItems
	if err := cur.All(ctx, &listItems); err != nil {
		return nil, err
	}
	total, err := db.collectionOrderItems().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	pagging.Total = total
	return listItems, nil
}
