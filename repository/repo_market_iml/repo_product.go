package repomarketiml

import (
	"context"
	"errors"
	"log"

	"github.com/KingSupermarket/controller/common"
	marketmodels "github.com/KingSupermarket/model/marketModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	db *mongo.Client
}

func NewDb(db *mongo.Client) *Db {
	return &Db{db: db}
}
func (db *Db) dbProducts() *mongo.Collection {
	return db.db.Database("KingSupermarket").Collection("products")
}
func (db *Db) CreateProduct(ctx context.Context, data *marketmodels.Product) error {
	count, err := db.db.Database("KingSupermarket").Collection("categories").
		CountDocuments(ctx, bson.M{"category_id": data.Category_id})
	if err != nil || count == 0 {
		return errors.New("category not found")
	}
	var cate *marketmodels.Category
	if err := db.db.Database("KingSupermarket").Collection("categories").
		FindOne(ctx, bson.M{"category_id": data.Category_id}).
		Decode(&cate); err != nil {
		return err
	}

	dataProduct, err := db.dbProducts().InsertOne(ctx, data)
	if err != nil {
		return err
	}
	oid, ok := dataProduct.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("faild oid")
	}
	data.Id = oid
	return nil
}
func (db *Db) GetProduct(ctx context.Context, id string) (*marketmodels.Product, error) {
	var product *marketmodels.Product
	if err := db.dbProducts().FindOne(ctx, bson.M{"product_id": id}).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}
func (db *Db) DeleteProduct(ctx context.Context, id string) error {
	del, err := db.dbProducts().DeleteOne(ctx, bson.M{"product_id": id})
	if err != nil {
		return err
	}
	if del.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}
func (db *Db) UpdateProduct(ctx context.Context, id string, data *marketmodels.Product) error {
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"title":       data.Title,
			"Image":       data.Image,
			"Description": data.Description,
			"price":       data.Price,
			"status":      data.Status,
			"updated_at":  data.Updated_at, // Added comma here
		}},
	}
	upd, err := db.dbProducts().UpdateOne(ctx, bson.M{"product_id": id}, update)
	if err != nil {
		return err
	}
	if upd.ModifiedCount == 0 {
		return errors.New("product not found or no changes made")
	}

	return nil
}
func (db *Db) GetListProduct(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Product, error) {
	pagging.Process()
	opts := options.Find().SetSkip(int64((pagging.Page - 1) * pagging.Limit)).SetLimit(int64(pagging.Limit))

	cur, err := db.dbProducts().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []*marketmodels.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}

	// Chuyển đổi từ []*marketmodels.Product sang []marketmodels.Product
	result := make([]marketmodels.Product, len(products))
	for i, product := range products {
		result[i] = *product // Chuyển từ *marketmodels.Product thành marketmodels.Product
	}

	total, err := db.dbProducts().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	pagging.Total = total
	return result, nil
}
