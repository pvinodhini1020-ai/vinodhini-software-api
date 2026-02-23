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

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id string) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id string) error
	List(page, pageSize int, search string, status string, clientID *string) ([]models.Project, int64, error)
	ListByEmployee(page, pageSize int, search string, status string, employeeID string) ([]models.Project, int64, error)
	AssignEmployees(projectID string, employeeIDs []string) error
}

type projectRepository struct {
	collection *mongo.Collection
	userColl   *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) ProjectRepository {
	return &projectRepository{
		collection: db.Collection("projects"),
		userColl:   db.Collection("users"),
	}
}

func (r *projectRepository) Create(project *models.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	if project.EmployeeIDs == nil {
		project.EmployeeIDs = []string{}
	}
	
	// Debug: Ensure project ID is set
	if project.ID == "" {
		return fmt.Errorf("project ID cannot be empty")
	}
	
	// Debug: Log before insertion
	fmt.Printf("Inserting project with _id: %s\n", project.ID)
	
	// Create document with explicit _id to ensure our custom ID is used
	doc := bson.M{
		"_id":         project.ID,
		"name":         project.Name,
		"description":  project.Description,
		"client_id":    project.ClientID,
		"status":       project.Status,
		"progress":     project.Progress,
		"employee_ids": project.EmployeeIDs,
		"created_at":   project.CreatedAt,
		"updated_at":   project.UpdatedAt,
	}
	
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}
	
	// Debug: Log the insertion result
	fmt.Printf("Insert result: %+v, InsertedID: %v\n", result, result.InsertedID)
	
	return nil
}

func (r *projectRepository) FindByID(id string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var project models.Project
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		return nil, err
	}

	var client models.User
	if err := r.userColl.FindOne(ctx, bson.M{"_id": project.ClientID}).Decode(&client); err == nil {
		project.Client = &client
	}

	if len(project.EmployeeIDs) > 0 {
		cursor, err := r.userColl.Find(ctx, bson.M{"_id": bson.M{"$in": project.EmployeeIDs}})
		if err == nil {
			cursor.All(ctx, &project.Employees)
		}
	}

	return &project, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	project.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": project.ID}, bson.M{"$set": project})
	return err
}

func (r *projectRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *projectRepository) List(page, pageSize int, search string, status string, clientID *string) ([]models.Project, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
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
	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, 0, err
	}

	for i := range projects {
		var client models.User
		if err := r.userColl.FindOne(ctx, bson.M{"_id": projects[i].ClientID}).Decode(&client); err == nil {
			projects[i].Client = &client
		}
	}

	return projects, total, nil
}


func (r *projectRepository) ListByEmployee(page, pageSize int, search string, status string, employeeID string) ([]models.Project, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"employee_ids": employeeID}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
			{"_id": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if status != "" {
		filter["status"] = status
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * pageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, 0, err
	}

	for i := range projects {
		var client models.User
		if err := r.userColl.FindOne(ctx, bson.M{"_id": projects[i].ClientID}).Decode(&client); err == nil {
			projects[i].Client = &client
		}
	}

	return projects, total, nil
}

func (r *projectRepository) AssignEmployees(projectID string, employeeIDs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": projectID},
		bson.M{"$set": bson.M{"employee_ids": employeeIDs, "updated_at": time.Now()}},
	)
	return err
}
