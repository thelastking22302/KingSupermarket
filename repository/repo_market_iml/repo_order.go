package repomarketiml

import (
	"context"
	"errors"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	usermodels "github.com/KingSupermarket/model/userModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDbOrder(db *mongo.Client) *Db {
	return &Db{db: db}
}

func (db *Db) collectionOrder() *mongo.Collection {
	return db.db.Database("KingSupermarket").Collection("order")
}

func (db *Db) CreateOrder(ctx context.Context, data *marketmodels.Order, idUser string) error {
	filter := bson.M{"user_id": idUser}
	cnt, err := db.db.Database("KingSupermarket").Collection("users").CountDocuments(ctx, filter)
	if err != nil || cnt == 0 {
		return err
	}
	var dataUser *usermodels.Users
	if err := db.db.Database("KingSupermarket").Collection("users").FindOne(ctx, filter).Decode(&dataUser); err != nil {
		return err
	}
	result, err := db.collectionOrder().InsertOne(ctx, &data)
	if err != nil {
		return err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("faild oid")
	}
	data.Id = oid
	return nil
}
func (db *Db) GetOrder(ctx context.Context, id string, idUser string) (*marketmodels.Order, error) {
	filter := bson.M{"user_id": idUser}
	cnt, err := db.db.Database("KingSupermarket").Collection("users").CountDocuments(ctx, filter)
	if err != nil || cnt == 0 {
		return nil, err
	}
	var dataUser *usermodels.Users
	if err := db.db.Database("KingSupermarket").Collection("users").FindOne(ctx, filter).Decode(&dataUser); err != nil {
		return nil, err
	}
	filterOrder := bson.M{"order_id": id}
	var dataOrder *marketmodels.Order
	if err := db.collectionOrder().FindOne(ctx, filterOrder).Decode(&dataOrder); err != nil {
		return nil, err
	}
	return dataOrder, nil
}
func (db *Db) GetListOrder(ctx context.Context, filter bson.M, pagging *common.Pagging, idUser string) ([]marketmodels.Order, error) {
	filterUser := bson.M{"user_id": idUser}
	cnt, err := db.db.Database("KingSupermarket").Collection("users").CountDocuments(ctx, filterUser)
	if err != nil || cnt == 0 {
		return nil, err
	}
	var dataUser *usermodels.Users
	if err := db.db.Database("KingSupermarket").Collection("users").FindOne(ctx, filterUser).Decode(&dataUser); err != nil {
		return nil, err
	}
	pagging.Process()
	opts := options.Find().SetSkip(int64((pagging.Page - 1) * pagging.Limit)).SetLimit(int64(pagging.Limit))
	cur, err := db.collectionOrder().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var listOrder []marketmodels.Order
	if err := cur.All(ctx, &listOrder); err != nil {
		return nil, err
	}
	total, err := db.collectionOrder().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	pagging.Total = total
	return listOrder, nil
}
func (db *Db) UpdateOrder(ctx context.Context, id string, data *marketmodels.Order, idUser string) error {
	filterUser := bson.M{"user_id": idUser}
	cnt, err := db.db.Database("KingSupermarket").Collection("users").CountDocuments(ctx, filterUser)
	if err != nil || cnt == 0 {
		return err
	}
	var dataUser *usermodels.Users
	if err := db.db.Database("KingSupermarket").Collection("users").FindOne(ctx, filterUser).Decode(&dataUser); err != nil {
		return err
	}
	upd := bson.D{
		{Key: "$set", Value: bson.M{
			"address":      data.Address,
			"phone_number": data.Phone_Number,
			"total_amount": data.Total_amount,
			"status":       data.Status,
			"notes":        data.Notes,
			"order_day":    data.Order_day,
			"updated_at":   data.Updated_at,
		}},
	}
	filter := bson.M{"order_id": id}
	dataUpd, err := db.collectionCategory().UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	if dataUpd.ModifiedCount == 0 {
		return errors.New("category not found or no changes made")
	}
	return nil
}
func (db *Db) DeleteOrder(ctx context.Context, id string, idUser string) error {
	filterUser := bson.M{"user_id": idUser}
	cnt, err := db.db.Database("KingSupermarket").Collection("users").CountDocuments(ctx, filterUser)
	if err != nil || cnt == 0 {
		return err
	}
	var dataUser *usermodels.Users
	if err := db.db.Database("KingSupermarket").Collection("users").FindOne(ctx, filterUser).Decode(&dataUser); err != nil {
		return err
	}
	filter := bson.M{"order_id": id}
	dataDel, err := db.collectionCategory().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if dataDel.DeletedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}
