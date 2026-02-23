package services

import (
	"errors"
	"fmt"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type ClientService interface {
	Create(req *models.CreateClientRequest) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(id string, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id string) error
	List(query *models.PaginationQuery) ([]models.User, int64, error)
}

type clientService struct {
	userRepo repositories.UserRepository
}

func NewClientService(userRepo repositories.UserRepository) ClientService {
	return &clientService{userRepo: userRepo}
}

func (s *clientService) Create(req *models.CreateClientRequest) (*models.User, error) {
	fmt.Printf("ClientService.Create called with Request: %+v\n", req)

	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Generate next user ID like USER01, USER02, etc.
	userID, err := s.userRepo.GetNextUserID()
	if err != nil {
		fmt.Printf("Error generating user ID: %v\n", err)
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Password hashing error: %v\n", err)
		return nil, errors.New("failed to hash password")
	}

	// Create user object
	user := &models.User{
		UserID:   userID,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Phone:     req.Phone,
		Role:      models.Role(req.Role),
		Company:   req.Company,
		Address:   req.Address,
		Status:    req.Status,
		Hide:      false, // Default to visible for new clients
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = models.RoleClient
	}

	fmt.Printf("Creating user: %+v\n", user)

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		fmt.Printf("Repository create error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Client created successfully\n")
	return user, nil
}

func (s *clientService) GetByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *clientService) Update(id string, req *models.UpdateUserRequest) (*models.User, error) {
	fmt.Printf("ClientService.Update called with ID: %s, Request: %+v\n", id, req)

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		fmt.Printf("Client not found: %v\n", err)
		return nil, errors.New("client not found")
	}

	fmt.Printf("Found client: %+v\n", user)

	if req.Name != "" {
		user.Name = req.Name
		fmt.Printf("Updated name to: %s\n", req.Name)
	}
	if req.Email != "" {
		user.Email = req.Email
		fmt.Printf("Updated email to: %s\n", req.Email)
	}
	if req.Phone != "" {
		user.Phone = req.Phone
		fmt.Printf("Updated phone to: %s\n", req.Phone)
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}
	if req.Company != "" {
		user.Company = req.Company
		fmt.Printf("Updated company to: %s\n", req.Company)
	}
	if req.Address != "" {
		user.Address = req.Address
		fmt.Printf("Updated address to: %s\n", req.Address)
	}
	if req.Salary > 0 {
		user.Salary = req.Salary
	}
	if req.Password != "" {
		// Hash password if provided
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Password hashing error: %v\n", err)
			return nil, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}
	if req.Status != "" {
		user.Status = req.Status
		fmt.Printf("Updated status to: %s\n", req.Status)
	}
	if req.Hide != nil {
		user.Hide = *req.Hide
		fmt.Printf("Updated hide to: %t\n", *req.Hide)
	}

	fmt.Printf("Updated client before save: %+v\n", user)

	if err := s.userRepo.Update(user); err != nil {
		fmt.Printf("Repository update error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Client updated successfully\n")
	return user, nil
}

func (s *clientService) Delete(id string) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("client not found")
	}

	return s.userRepo.Delete(id)
}

func (s *clientService) List(query *models.PaginationQuery) ([]models.User, int64, error) {
	return s.userRepo.List(query.Page, query.PageSize, query.Search, "client")
}
