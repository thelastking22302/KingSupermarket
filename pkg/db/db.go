package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectionMongo(userName, Password string) (*mongo.Client, error) {
	conStr := fmt.Sprintf("mongodb://%s:%s@localhost:27017/", userName, Password)
	client, err := mongo.NewClient(options.Client().ApplyURI(conStr))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client, nil
}
