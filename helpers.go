package mongodb

import (
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GenerateMongoID generates a new MongoDB ObjectID.
func GenerateMongoID() bson.ObjectID {
	return bson.NewObjectID()
}

// GenerateNanoID generates a new NanoID.
func GenerateNanoID() string {
	id, err := gonanoid.New()
	if err != nil {
		log.Fatalf("Failed to generate NanoID: %v", err)
	}
	return id
}
