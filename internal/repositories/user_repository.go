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

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByUserID(userID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	List(page, pageSize int, search string, role string) ([]models.User, int64, error)
	GetNextUserID() (string, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{collection: db.Collection("users")}
}

func (r *userRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	// Set the UserID as the MongoDB _id
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *userRepository) FindByUserID(userID string) (*models.User, error) {
	// Since UserID is now the _id, we can use FindByID
	return r.FindByID(userID)
}

func (r *userRepository) GetNextUserID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find the highest existing user_id
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(1)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	var lastUser models.User
	if cursor.Next(ctx) {
		if err := cursor.Decode(&lastUser); err != nil {
			return "", err
		}
	}

	// If no users exist, start with USER01
	if lastUser.UserID == "" {
		return "USER01", nil
	}

	// Extract the numeric part and increment
	var lastNum int
	if _, err := fmt.Sscanf(lastUser.UserID, "USER%d", &lastNum); err != nil {
		// If parsing fails, start with USER01
		return "USER01", nil
	}

	nextNum := lastNum + 1
	return fmt.Sprintf("USER%02d", nextNum), nil
}

func (r *userRepository) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()
	
	// Create update document excluding UserID (_id) field
	updateDoc := bson.M{
		"email":      user.Email,
		"password":   user.Password,
		"name":       user.Name,
		"phone":      user.Phone,
		"role":       user.Role,
		"department": user.Department,
		"company":    user.Company,
		"address":    user.Address,
		"salary":     user.Salary,
		"status":     user.Status,
		"hide":       user.Hide,
		"updated_at": user.UpdatedAt,
	}
	
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.UserID}, bson.M{"$set": updateDoc})
	return err
}

func (r *userRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *userRepository) List(page, pageSize int, search string, role string) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"company": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if role != "" {
		filter["role"] = role
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

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
