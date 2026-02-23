package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CounterRepository interface {
	GetNextSequence(counterName string) (int, error)
}

type counterRepository struct {
	collection *mongo.Collection
}

func NewCounterRepository(db *mongo.Database) CounterRepository {
	return &counterRepository{
		collection: db.Collection("counters"),
	}
}

func (r *counterRepository) GetNextSequence(counterName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use findOneAndUpdate for atomic sequence generation
	// This ensures concurrency safety by using MongoDB's atomic operations
	filter := bson.M{"_id": counterName}
	update := bson.M{"$inc": bson.M{"sequence": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		ID       string `bson:"_id"`
		Sequence int    `bson:"sequence"`
	}

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// First time creating, start with sequence 1
			_, err := r.collection.InsertOne(ctx, bson.M{"_id": counterName, "sequence": 1})
			if err != nil {
				return 1, fmt.Errorf("failed to initialize %s counter: %w", counterName, err)
			}
			return 1, nil
		}
		return 1, fmt.Errorf("failed to get next %s sequence: %w", counterName, err)
	}

	return result.Sequence, nil
}
