package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"log"
)

func Connect(databaseName, collectionName string) (*mongo.Client, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB Connection Error:", err)
	}
	collection := client.Database(databaseName).Collection(collectionName)
	return client, collection
}