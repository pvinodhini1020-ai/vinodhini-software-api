package services

import (
	"errors"
	"fmt"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(id string) (*models.User, error)
	Update(id string, req *models.UpdateUserRequest) (*models.User, error)
	Delete(id string) error
	List(query *models.PaginationQuery, role string) ([]models.User, int64, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) Update(id string, req *models.UpdateUserRequest) (*models.User, error) {
	fmt.Printf("UserService.Update called with ID: %s, Request: %+v\n", id, req)
	
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		fmt.Printf("User not found: %v\n", err)
		return nil, errors.New("user not found")
	}
	
	fmt.Printf("Found user: %+v\n", user)

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
		fmt.Printf("Updated phone to: %s\n", req.Phone)
	}
	if req.Role != "" {
		user.Role = models.Role(req.Role)
	}
	if req.Department != "" {
		user.Department = req.Department
		fmt.Printf("Updated department to: %s\n", req.Department)
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
	}
	
	fmt.Printf("Updated user before save: %+v\n", user)

	if err := s.userRepo.Update(user); err != nil {
		fmt.Printf("Repository update error: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("User updated successfully\n")

	return user, nil
}

func (s *userService) Delete(id string) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) List(query *models.PaginationQuery, role string) ([]models.User, int64, error) {
	return s.userRepo.List(query.Page, query.PageSize, query.Search, role)
}
