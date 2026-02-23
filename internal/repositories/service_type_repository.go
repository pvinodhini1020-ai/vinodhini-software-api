package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/vinodhini/software-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceTypeRepository struct {
	collection *mongo.Collection
}

func NewServiceTypeRepository(db *mongo.Database) *ServiceTypeRepository {
	return &ServiceTypeRepository{
		collection: db.Collection("service_types"),
	}
}

func (r *ServiceTypeRepository) Create(serviceType *models.ServiceType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceType.ID = primitive.NewObjectID().Hex()
	serviceType.CreatedAt = time.Now()
	serviceType.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, serviceType)
	return err
}

func (r *ServiceTypeRepository) GetByID(id string) (*models.ServiceType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var serviceType models.ServiceType
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&serviceType)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("service type not found")
		}
		return nil, err
	}

	return &serviceType, nil
}

func (r *ServiceTypeRepository) GetAll(status *string) ([]models.ServiceType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if status != nil {
		filter["status"] = *status
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var serviceTypes []models.ServiceType
	if err = cursor.All(ctx, &serviceTypes); err != nil {
		return nil, err
	}

	return serviceTypes, nil
}

func (r *ServiceTypeRepository) Update(id string, serviceType *models.ServiceType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceType.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":        serviceType.Name,
			"description": serviceType.Description,
			"status":      serviceType.Status,
			"updated_at":  serviceType.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("service type not found")
	}

	return nil
}

func (r *ServiceTypeRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("service type not found")
	}

	return nil
}
