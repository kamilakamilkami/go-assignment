package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect function to initialize and return a MongoDB client
func Connect() *mongo.Client {
	// MongoDB URI (your MongoDB server URL)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Change if using remote DB

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB Ping Error:", err)
	}
	log.Println("Successfully connected to MongoDB!")

	return client
}
