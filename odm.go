package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ODM[T any] struct {
	Collection *mongo.Collection
	CreateOne  func(props T) (*mongo.InsertOneResult, error)
	FindOne    func(filter bson.M) (*T, error)
	FindMany   func(filter bson.M) ([]T, error)
	FindById   func(id string) (*T, error)
	List       func() ([]T, error)
	UpdateOne  func(filter bson.M, props bson.M) (*mongo.UpdateResult, error)
	UpdateMany func(filter bson.M, props bson.M) (*mongo.UpdateResult, error)
	UpdateById func(id string, props bson.M) (*mongo.UpdateResult, error)
	DeleteOne  func(filter bson.M) (*mongo.DeleteResult, error)
	DeleteMany func(filter bson.M) (*mongo.DeleteResult, error)
	DeleteById func(id string) (*mongo.DeleteResult, error)
}

func NewODM[T any](collection *mongo.Collection) *ODM[T] {
	return &ODM[T]{
		Collection: collection,
		CreateOne:  createOneGenerator[T](collection),
		FindOne:    findOneGenerator[T](collection),
		FindMany:   findManyGenerator[T](collection),
		FindById:   findByIdGenerator[T](collection),
		List:       listGenerator[T](collection),
		UpdateOne:  updateOneGenerator[T](collection),
		UpdateMany: updateManyGenerator[T](collection),
		UpdateById: updateByIdGenerator[T](collection),
		DeleteOne:  deleteOneGenerator[T](collection),
		DeleteMany: deleteManyGenerator[T](collection),
		DeleteById: deleteByIdGenerator[T](collection),
	}
}

func createOneGenerator[T any](collection *mongo.Collection) func(props T) (*mongo.InsertOneResult, error) {
	return func(props T) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(context.Background(), props)
	}
}

func findOneGenerator[T any](collection *mongo.Collection) func(filter bson.M) (*T, error) {
	return func(filter bson.M) (*T, error) {
		var doc T
		err := collection.FindOne(context.Background(), filter).Decode(&doc)
		if err != nil {
			return nil, err
		}
		return &doc, nil
	}
}

func findManyGenerator[T any](collection *mongo.Collection) func(filter bson.M) ([]T, error) {
	return func(filter bson.M) ([]T, error) {
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
}

func findByIdGenerator[T any](collection *mongo.Collection) func(id string) (*T, error) {
	return func(id string) (*T, error) {
		var doc T
		err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
		if err != nil {
			return nil, err
		}
		return &doc, nil
	}
}

func listGenerator[T any](collection *mongo.Collection) func() ([]T, error) {
	return func() ([]T, error) {
		cursor, err := collection.Find(context.Background(), bson.M{})
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
}

func updateOneGenerator[T any](collection *mongo.Collection) func(filter bson.M, props bson.M) (*mongo.UpdateResult, error) {
	return func(filter bson.M, props bson.M) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(context.Background(), filter, props)
	}
}

func updateManyGenerator[T any](collection *mongo.Collection) func(filter bson.M, props bson.M) (*mongo.UpdateResult, error) {
	return func(filter bson.M, props bson.M) (*mongo.UpdateResult, error) {
		return collection.UpdateMany(context.Background(), filter, props)
	}
}

func updateByIdGenerator[T any](collection *mongo.Collection) func(id string, props bson.M) (*mongo.UpdateResult, error) {
	return func(id string, props bson.M) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(context.Background(), bson.M{"_id": id}, props)
	}
}

func deleteOneGenerator[T any](collection *mongo.Collection) func(filter bson.M) (*mongo.DeleteResult, error) {
	return func(filter bson.M) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(context.Background(), filter)
	}
}

func deleteManyGenerator[T any](collection *mongo.Collection) func(filter bson.M) (*mongo.DeleteResult, error) {
	return func(filter bson.M) (*mongo.DeleteResult, error) {
		return collection.DeleteMany(context.Background(), filter)
	}
}

func deleteByIdGenerator[T any](collection *mongo.Collection) func(id string) (*mongo.DeleteResult, error) {
	return func(id string) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(context.Background(), bson.M{"_id": id})
	}
}
