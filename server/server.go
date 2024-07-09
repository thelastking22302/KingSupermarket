package server

import (
	"fmt"
	"log"
	"os"

	"github.com/KingSupermarket/pkg/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewConnectionMongo() *mongo.Client {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	userName := os.Getenv("USER_NAME_MONGO")
	password := os.Getenv("PASSWORD_MONGO")
	clientMongo, err := db.ConnectionMongo(userName, password)
	if err != nil {
		log.Fatalf("Fails to connect to Mongo: %v", err)
	}
	fmt.Println("Connect Mongo success")
	return clientMongo
}
