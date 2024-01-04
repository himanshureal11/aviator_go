package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDB *mongo.Database
)

// InitMongoDB initializes the MongoDB connection.
func InitMongoDB() {
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		log.Println("MONGODB_URL environment variable is not set.")
	}

	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Println("Error creating MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pinging MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Set the global MongoDB database variable
	MongoDB = client.Database("real11-test") // Replace with your actual database name

}

// GetMongoDB returns the MongoDB client.
func GetMongoDB() *mongo.Database {
	if MongoDB == nil {
		InitMongoDB()
	}
	return MongoDB
}

func CreateCollection(collectionName string) *mongo.Collection {
	// Use the GetMongoDB function from the package where it's defined
	db := GetMongoDB()

	// Create a collection with the specified name
	collection := db.Collection(collectionName)

	// Return the reference to the collection
	return collection
}
