package mongodb_odm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateOne[T any](collection *mongo.Collection, props T, getNew bool) (*T, error) {
	_, err := collection.InsertOne(context.Background(), props)
	if err != nil {
		return nil, err
	}

	if !getNew {
		return nil, nil
	}

	doc, err := FindById[T](collection, "id")
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func FindOne[T any](collection *mongo.Collection, filter bson.M) (*T, error) {
	var doc T
	err := collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func FindMany[T any](collection *mongo.Collection, filter bson.M) ([]T, error) {
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var docs []T
	if err := cursor.All(context.Background(), &docs); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return []T{}, nil
	}
	return docs, nil
}

func FindById[T any](collection *mongo.Collection, id string) (*T, error) {
	doc, err := FindOne[T](collection, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func List[T any](collection *mongo.Collection) ([]T, error) {
	return FindMany[T](collection, bson.M{})
}

func UpdateOne[T any](collection *mongo.Collection, filter bson.M, props bson.M, getNew bool) (*T, error) {
	_, err := collection.UpdateOne(context.Background(), filter, props)
	if err != nil {
		return nil, err
	}
	if !getNew {
		return nil, nil
	}

	doc, err := FindOne[T](collection, filter)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func UpdateMany[T any](collection *mongo.Collection, filter bson.M, props bson.M, getNew bool) ([]T, error) {
	_, err := collection.UpdateMany(context.Background(), filter, props)
	if err != nil {
		return nil, err
	}
	if !getNew {
		return nil, nil
	}

	docs, err := FindMany[T](collection, filter)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func UpdateById[T any](collection *mongo.Collection, id string, props bson.M, getNew bool) (*T, error) {
	doc, err := UpdateOne[T](collection, bson.M{"_id": id}, props, getNew)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func DeleteOne[T any](collection *mongo.Collection, filter bson.M) error {
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMany[T any](collection *mongo.Collection, filter bson.M) error {
	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteById[T any](collection *mongo.Collection, id string) error {
	err := DeleteOne[T](collection, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
