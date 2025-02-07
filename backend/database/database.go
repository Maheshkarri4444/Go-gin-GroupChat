package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = createMonogClient()

func createMonogClient() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading the file")
	}
	MongodbUri := os.Getenv("MONGODB_URI")

	if MongodbUri == "" {
		log.Fatal("Mongodburi is not found in the env variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongodbUri))
	if err != nil {
		log.Fatal("mongo db connection error: ", err)
	}

	fmt.Println("connected to mongodb")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("group-chat").Collection(collectionName)

}
