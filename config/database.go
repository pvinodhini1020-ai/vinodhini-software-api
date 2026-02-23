package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(cfg *Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI).SetMaxPoolSize(100)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("MongoDB connected successfully")
	return client.Database(cfg.MongoDB.Database), nil
}

func CreateIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("service_requests").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
	})

	log.Println("Indexes created successfully")
	return err
}
