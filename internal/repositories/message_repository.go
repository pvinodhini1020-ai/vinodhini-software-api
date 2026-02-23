package repositories

import (
	"context"
	"time"

	"github.com/vinodhini/software-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	Create(message *models.Message) error
	FindByID(id string) (*models.Message, error)
	Delete(id string) error
	ListByProject(projectID string, page, pageSize int) ([]models.Message, int64, error)
}

type messageRepository struct {
	collection *mongo.Collection
	userColl   *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"),
		userColl:   db.Collection("users"),
	}
}

func (r *messageRepository) Create(message *models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	
	// Create document with explicit _id to ensure our custom ID is used
	doc := bson.M{
		"_id":         message.ID,
		"content":      message.Content,
		"sender_id":    message.SenderID,
		"project_id":   message.ProjectID,
		"created_at":   message.CreatedAt,
		"updated_at":   message.UpdatedAt,
	}
	
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *messageRepository) FindByID(id string) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var message models.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}

	var sender models.User
	if err := r.userColl.FindOne(ctx, bson.M{"_id": message.SenderID}).Decode(&sender); err == nil {
		message.Sender = &sender
	}

	return &message, nil
}

func (r *messageRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}


func (r *messageRepository) ListByProject(projectID string, page, pageSize int) ([]models.Message, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"project_id": projectID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	for i := range messages {
		var sender models.User
		if err := r.userColl.FindOne(ctx, bson.M{"_id": messages[i].SenderID}).Decode(&sender); err == nil {
			messages[i].Sender = &sender
		}
	}

	return messages, total, nil
}
