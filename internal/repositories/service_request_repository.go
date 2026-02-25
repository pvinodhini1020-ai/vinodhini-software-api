package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/vinodhini/software-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceRequestRepository interface {
	Create(request *models.ServiceRequest) error
	FindByID(id string) (*models.ServiceRequest, error)
	Update(request *models.ServiceRequest) error
	Delete(id string) error
	List(page, pageSize int, search string, status string, clientID *string) ([]models.ServiceRequest, int64, error)
}

type serviceRequestRepository struct {
	collection  *mongo.Collection
	userColl    *mongo.Collection
	projectColl *mongo.Collection
}

func NewServiceRequestRepository(db *mongo.Database) ServiceRequestRepository {
	return &serviceRequestRepository{
		collection:  db.Collection("service_requests"),
		userColl:    db.Collection("users"),
		projectColl: db.Collection("projects"),
	}
}

func (r *serviceRequestRepository) Create(request *models.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	
	// Create document with explicit _id to ensure our custom ID is used
	doc := bson.M{
		"_id":         request.ID,
		"title":        request.Title,
		"description":  request.Description,
		"client_id":    request.ClientID,
		"project_id":   request.ProjectID,
		"status":       request.Status,
		"created_at":   request.CreatedAt,
		"updated_at":   request.UpdatedAt,
	}
	
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *serviceRequestRepository) FindByID(id string) (*models.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request models.ServiceRequest
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err != nil {
		return nil, err
	}

	var client models.User
	clientErr := r.userColl.FindOne(ctx, bson.M{"_id": request.ClientID}).Decode(&client)
	if clientErr == nil {
		request.Client = &client
		fmt.Printf("Successfully loaded client %s for service request %s\n", client.Name, request.ID)
	} else {
		fmt.Printf("Failed to find client for service request %s, client ID %s: %v\n", request.ID, request.ClientID, clientErr)
		// Create a default client object to prevent null reference issues
		request.Client = &models.User{
			UserID: request.ClientID,
			Name:   "Unknown Client",
			Email:  "unknown@example.com",
			Role:   models.RoleClient,
		}
	}

	if request.ProjectID != nil {
		var project models.Project
		if err := r.projectColl.FindOne(ctx, bson.M{"_id": *request.ProjectID}).Decode(&project); err == nil {
			request.Project = &project
		}
	}

	return &request, nil
}

func (r *serviceRequestRepository) Update(request *models.ServiceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": request.ID}, bson.M{"$set": request})
	return err
}

func (r *serviceRequestRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}


func (r *serviceRequestRepository) List(page, pageSize int, search string, status string, clientID *string) ([]models.ServiceRequest, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
			{"_id": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if status != "" {
		filter["status"] = status
	}
	if clientID != nil {
		filter["client_id"] = *clientID
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize)).SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var requests []models.ServiceRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, 0, err
	}

	for i := range requests {
		var client models.User
		clientErr := r.userColl.FindOne(ctx, bson.M{"_id": requests[i].ClientID}).Decode(&client)
		if clientErr == nil {
			requests[i].Client = &client
		} else {
			fmt.Printf("Failed to find client for service request %s, client ID %s: %v\n", requests[i].ID, requests[i].ClientID, clientErr)
			// Create a default client object to prevent null reference issues
			requests[i].Client = &models.User{
				UserID: requests[i].ClientID,
				Name:   "Unknown Client",
				Email:  "unknown@example.com",
				Role:   models.RoleClient,
			}
		}
	}

	return requests, total, nil
}
