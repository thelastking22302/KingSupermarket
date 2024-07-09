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

func NewDbInvoice(db *mongo.Client) *Db {
	return &Db{db: db}
}

func (db *Db) collectionInvoice() *mongo.Collection {
	return db.db.Database("KingSupermarket").Collection("invoice")
}

func (db *Db) CreateInvoice(ctx context.Context, data *marketmodels.Invoice) error {
	filterOrder := bson.M{"order_id": data.Order_id}
	count, err := db.collectionOrder().CountDocuments(ctx, filterOrder)
	if err != nil && count == 0 {
		return errors.New("invoice not found")
	}
	var orderData marketmodels.Order
	if err := db.collectionOrder().FindOne(ctx, filterOrder).Decode(&orderData); err != nil {
		return errors.New("order not found")
	}
	indata, err := db.collectionInvoice().InsertOne(ctx, &data)
	if err != nil {
		return errors.New("insert failer")
	}
	oid, ok := indata.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("insert objectID failer")
	}
	data.ID = oid
	return nil
}
func (db *Db) GetInvoice(ctx context.Context, id string) (*marketmodels.Invoice, error) {
	var dataInvoice marketmodels.Invoice
	filter := bson.M{"invoice_id": id}
	if err := db.collectionInvoice().FindOne(ctx, filter).Decode(&dataInvoice); err != nil {
		return nil, err
	}
	return &dataInvoice, nil
}
func (db *Db) GetListInvoice(ctx context.Context, filter bson.M, pagging *common.Pagging) ([]marketmodels.Invoice, error) {
	pagging.Process()
	opts := options.Find().SetSkip((int64((pagging.Page - 1) - pagging.Limit))).SetLimit(int64(pagging.Limit))

	cur, err := db.collectionInvoice().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var listInvoice []marketmodels.Invoice
	if err := cur.All(ctx, &listInvoice); err != nil {
		return nil, err
	}
	total, err := db.collectionInvoice().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	pagging.Total = total
	return listInvoice, nil
}
func (db *Db) UpdateInvoice(ctx context.Context, id string, data *marketmodels.Invoice) error {
	upd := bson.D{
		{
			Key: "$set", Value: bson.M{
				"payment_method":   data.Payment_method,
				"payment_status":   data.Payment_status,
				"Payment_due_date": data.Payment_due_date,
				"updated_at":       data.Updated_at,
			},
		},
	}
	filter := bson.M{"invoice_id": id}
	count, err := db.collectionInvoice().UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	if count.ModifiedCount == 0 {
		return errors.New("invoice not found or no changes made")
	}
	return nil
}
func (db *Db) DeleteInvoice(ctx context.Context, id string) error {
	filter := bson.M{"invoice_id": id}
	cnt, err := db.collectionInvoice().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if cnt.DeletedCount == 0 {
		return errors.New("invoice not found")
	}
	return nil
}
