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

func NewDBCategory(db *mongo.Client) *Db {
	return &Db{db: db}
}

func (db *Db) collectionCategory() *mongo.Collection {
	return db.db.Database("KingSupermarket").Collection("categories")
}

func (db *Db) CreateCategory(ctx context.Context, data *marketmodels.Category) error {
	dataCategory, err := db.collectionCategory().InsertOne(ctx, &data)
	if err != nil {
		return err
	}
	if oid, ok := dataCategory.InsertedID.(primitive.ObjectID); ok {
		data.Id = oid
	} else {
		return errors.New("failed to convert InsertedID to primitive.ObjectID")
	}
	return nil
}

func (db *Db) GetCategory(ctx context.Context, id string) (*marketmodels.Category, error) {
	var data *marketmodels.Category
	filter := bson.M{"category_id": id}
	if err := db.collectionCategory().FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (db *Db) GetListCategory(ctx context.Context, pagging *common.Pagging) ([]marketmodels.Category, error) {
	pagging.Process()
	opts := options.Find().SetSkip((int64((pagging.Page - 1) * pagging.Limit))).SetLimit(int64(pagging.Limit))
	cur, err := db.collectionCategory().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var listCate []marketmodels.Category
	if err := cur.All(ctx, &listCate); err != nil {
		return nil, err
	}
	total, err := db.collectionCategory().CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	pagging.Total = total
	return listCate, nil
}

func (db *Db) UpdateCategory(ctx context.Context, id string, data *marketmodels.Category) error {
	upd := bson.D{
		{Key: "$set", Value: bson.M{
			"name":       data.Name,
			"updated_at": data.Updated_at,
		}},
	}
	filter := bson.M{"category_id": id}
	dataUpd, err := db.collectionCategory().UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	if dataUpd.ModifiedCount == 0 {
		return errors.New("category not found or no changes made")
	}
	return nil
}
func (db *Db) DeleteCategory(ctx context.Context, id string) error {
	filter := bson.M{"category_id": id}
	dataDel, err := db.collectionCategory().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if dataDel.DeletedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}
