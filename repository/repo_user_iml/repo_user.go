package repouseriml

import (
	"context"
	"errors"
	"log"
	"time"

	marketmodels "github.com/KingSupermarket/model/marketModels"
	usermodels "github.com/KingSupermarket/model/userModels"
	requsermodel "github.com/KingSupermarket/model/userModels/reqUserModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	mg *mongo.Client
}

func NewDB(mg *mongo.Client) *db {
	return &db{mg: mg}
}

func (dbm *db) HistoryPurchases(ctx context.Context, id string) (*usermodels.Users, []marketmodels.OrderItems, error) {
	var dataUsers *usermodels.Users
	var dataProducts []marketmodels.OrderItems

	// Tìm kiếm người dùng
	userFilter := bson.M{"user_id": id}
	err := dbm.mg.Database("KingSupermarket").Collection("users").FindOne(ctx, userFilter).Decode(&dataUsers)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, nil // Không tìm thấy người dùng
		}
		return nil, nil, err
	}

	// Tìm kiếm các sản phẩm đã mua
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": id}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "order",
			"localField":   "user_id",
			"foreignField": "user_id",
			"as":           "order",
		}}},
		{{Key: "$unwind", Value: "$order"}},
		{{Key: "$unwind", Value: "$order.order_item"}},
		{{Key: "$replaceRoot", Value: bson.M{
			"newRoot": bson.M{
				"user_id":    "$user_id",
				"product_id": "$order.order_item.product_id",
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{"user_id": "$user_id", "product_id": "$product_id"},
		}}},
	}

	cursor, err := dbm.mg.Database("KingSupermarket").Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item marketmodels.OrderItems
		if err := cursor.Decode(&item); err != nil {
			return nil, nil, err
		}
		dataProducts = append(dataProducts, item)
	}

	return dataUsers, dataProducts, nil
}
func (dbm *db) SignIn(ctx context.Context, data *requsermodel.SigninModel) (*usermodels.Users, error) {
	var foundData usermodels.Users
	if err := dbm.mg.Database("KingSupermarket").Collection("users").FindOne(ctx, bson.M{"email": data.Email}).Decode(&foundData); err != nil {
		return nil, err
	}

	return &foundData, nil
}

func (dbm *db) SignUp(ctx context.Context, data *usermodels.Users) (*usermodels.Users, error) {
	countEmail, err := dbm.mg.Database("KingSupermarket").Collection("users").CountDocuments(ctx, bson.M{"email": data.Email})
	if err != nil {
		log.Fatalln("error occured while checking for the email", err)
	}
	if countEmail > 0 {
		log.Fatalln("this email already exsits")
	}

	result, insertErr := dbm.mg.Database("KingSupermarket").Collection("users").InsertOne(ctx, &data)
	if insertErr != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		data.Id = oid
	} else {
		return nil, errors.New("failed to convert InsertedID to primitive.ObjectID")
	}
	return data, nil
}
func (dbm *db) ProfileUser(ctx context.Context, id string) (*usermodels.Users, error) {
	var dataUser *usermodels.Users
	if err := dbm.mg.Database("KingSupermarket").Collection("users").FindOne(ctx, bson.M{"user_id": id}).Decode(&dataUser); err != nil {
		return nil, err
	}
	return dataUser, nil
}
func (dbm *db) UpdateUser(ctx context.Context, id string, data *usermodels.Users) error {
	var updateObj primitive.D
	if data.Name == nil {
		updateObj = append(updateObj, bson.E{Key: "name", Value: data.Name})
	}
	if data.Email == nil {
		updateObj = append(updateObj, bson.E{Key: "email", Value: data.Email})
	}
	if data.Password == "" {
		updateObj = append(updateObj, bson.E{Key: "password", Value: data.Password})
	}
	if data.Avatar == nil {
		updateObj = append(updateObj, bson.E{Key: "avatar", Value: data.Avatar})
	}
	if data.Phone_Number == nil {
		updateObj = append(updateObj, bson.E{Key: "phone_number", Value: data.Phone_Number})
	}
	timeUpdate := time.Now().UTC()
	data.Updated_At = &timeUpdate
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: data.Updated_At})
	_, err := dbm.mg.Database("KingSupermarket").Collection("users").UpdateOne(ctx, bson.M{"user_id": id}, bson.D{{Key: "$set", Value: updateObj}})
	if err != nil {
		return err
	}
	return nil
}

func (dbm *db) DeleteUser(ctx context.Context, id string) error {
	_, err := dbm.mg.Database("KingSupermarket").Collection("users").DeleteOne(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}
	return nil
}
