package mongodb

import (
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var defaultAlphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateMongoID() bson.ObjectID {
	return bson.NewObjectID()
}

func GenerateNanoID() string {
	return GenerateNanoIdBySize(24)
}

func GenerateNanoIdBySize(size int) string {
	return GenerateCustomNanoID(defaultAlphabets, size)
}

func GenerateCustomNanoID(alphabets string, size int) string {
	id, err := gonanoid.Generate(alphabets, size)
	if err != nil {
		log.Printf("Error generating NanoID: %v", err)
		return ""
	}
	return id
}
