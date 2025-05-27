package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

// Connect initializes the MongoDB client using the MONGODB_URI environment variable.
func Connect() error {
	log.Println("[MongoDB] Connecting to MongoDB...")
	if Config.URI == "" {
		log.Fatal("[MongoDB] Connection failed. Invalid configuration.")
		return nil
	}

	var err error
	Client, err = mongo.Connect(options.Client().ApplyURI(Config.URI))
	if err != nil {
		log.Fatal("[MongoDB] Connection failed. " + err.Error())
		return err
	}
	Database = Client.Database(Config.Database)
	log.Println("[MongoDB] Connection successful.")
	return nil
}

// Disconnect closes the connection to the MongoDB client.
func Disconnect() error {
	if Client == nil {
		return nil
	}
	return Client.Disconnect(context.TODO())
}
