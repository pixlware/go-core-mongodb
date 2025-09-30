package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ODM[T any] struct {
	Collection *mongo.Collection
	CreateOne  func(doc T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	CreateMany func(docs []T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error)
	BulkWrite  func(models []mongo.WriteModel) (*mongo.BulkWriteResult, error)
	FindOne    func(filter bson.M, opts ...options.Lister[options.FindOneOptions]) (*T, error)
	FindMany   func(filter bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error)
	FindById   func(id string, opts ...options.Lister[options.FindOneOptions]) (*T, error)
	List       func(opts ...options.Lister[options.FindOptions]) ([]T, error)
	Count      func(filter bson.M, opts ...options.Lister[options.CountOptions]) (int64, error)
	UpdateOne  func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	UpdateMany func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error)
	UpdateById func(id string, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	DeleteOne  func(filter bson.M, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
	DeleteMany func(filter bson.M, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error)
	DeleteById func(id string, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
}

func NewODM[T any](collection *mongo.Collection) *ODM[T] {
	return &ODM[T]{
		Collection: collection,
		CreateOne:  createOneGenerator[T](collection),
		CreateMany: createManyGenerator[T](collection),
		BulkWrite:  bulkWriteGenerator[T](collection),
		FindOne:    findOneGenerator[T](collection),
		FindMany:   findManyGenerator[T](collection),
		FindById:   findByIdGenerator[T](collection),
		List:       listGenerator[T](collection),
		Count:      countGenerator[T](collection),
		UpdateOne:  updateOneGenerator[T](collection),
		UpdateMany: updateManyGenerator[T](collection),
		UpdateById: updateByIdGenerator[T](collection),
		DeleteOne:  deleteOneGenerator[T](collection),
		DeleteMany: deleteManyGenerator[T](collection),
		DeleteById: deleteByIdGenerator[T](collection),
	}
}

func createOneGenerator[T any](collection *mongo.Collection) func(doc T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
	return func(doc T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(context.Background(), doc, opts...)
	}
}

func createManyGenerator[T any](collection *mongo.Collection) func(docs []T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
	return func(docs []T, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
		return collection.InsertMany(context.Background(), docs, opts...)
	}
}

func bulkWriteGenerator[T any](collection *mongo.Collection) func(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return func(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
		bulkOptions := options.BulkWrite().SetOrdered(true)
		return collection.BulkWrite(context.Background(), models, bulkOptions)
	}
}

func findOneGenerator[T any](collection *mongo.Collection) func(filter bson.M, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	return func(filter bson.M, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
		var doc T
		err := collection.FindOne(context.Background(), filter, opts...).Decode(&doc)
		if err != nil {
			return nil, err
		}
		return &doc, nil
	}
}

func findManyGenerator[T any](collection *mongo.Collection) func(filter bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	return func(filter bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error) {
		cursor, err := collection.Find(context.Background(), filter, opts...)
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

func findByIdGenerator[T any](collection *mongo.Collection) func(id string, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
	return func(id string, opts ...options.Lister[options.FindOneOptions]) (*T, error) {
		var doc T
		err := collection.FindOne(context.Background(), bson.M{"_id": id}, opts...).Decode(&doc)
		if err != nil {
			return nil, err
		}
		return &doc, nil
	}
}

func listGenerator[T any](collection *mongo.Collection) func(opts ...options.Lister[options.FindOptions]) ([]T, error) {
	return func(opts ...options.Lister[options.FindOptions]) ([]T, error) {
		cursor, err := collection.Find(context.Background(), bson.M{}, opts...)
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

func countGenerator[T any](collection *mongo.Collection) func(filter bson.M, opts ...options.Lister[options.CountOptions]) (int64, error) {
	return func(filter bson.M, opts ...options.Lister[options.CountOptions]) (int64, error) {
		return collection.CountDocuments(context.Background(), filter, opts...)
	}
}

func updateOneGenerator[T any](collection *mongo.Collection) func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(context.Background(), filter, update, opts...)
	}
}

func updateManyGenerator[T any](collection *mongo.Collection) func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
	return func(filter bson.M, update bson.M, opts ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
		return collection.UpdateMany(context.Background(), filter, update, opts...)
	}
}

func updateByIdGenerator[T any](collection *mongo.Collection) func(id string, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	return func(id string, update bson.M, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(context.Background(), bson.M{"_id": id}, update, opts...)
	}
}

func deleteOneGenerator[T any](collection *mongo.Collection) func(filter bson.M, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return func(filter bson.M, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(context.Background(), filter, opts...)
	}
}

func deleteManyGenerator[T any](collection *mongo.Collection) func(filter bson.M, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error) {
	return func(filter bson.M, opts ...options.Lister[options.DeleteManyOptions]) (*mongo.DeleteResult, error) {
		return collection.DeleteMany(context.Background(), filter, opts...)
	}
}

func deleteByIdGenerator[T any](collection *mongo.Collection) func(id string, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
	return func(id string, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(context.Background(), bson.M{"_id": id}, opts...)
	}
}
