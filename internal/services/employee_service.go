package services

import (
	"errors"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService interface {
	Create(req *models.CreateEmployeeRequest) (*models.User, error)
	GetByID(id string) (*models.User, error)
	List(query *models.PaginationQuery) ([]models.User, int64, error)
}

type employeeService struct {
	employeeRepo repositories.EmployeeRepository
	userRepo     repositories.UserRepository
}

func NewEmployeeService(employeeRepo repositories.EmployeeRepository, userRepo repositories.UserRepository) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
	}
}

func (s *employeeService) Create(req *models.CreateEmployeeRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Get next user ID
	userID, err := s.userRepo.GetNextUserID()
	if err != nil {
		return nil, errors.New("failed to generate user ID")
	}

	// Create employee
	employee := &models.User{
		UserID:   userID,
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     models.RoleEmployee,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, errors.New("failed to create employee")
	}

	// Clear password before returning
	employee.Password = ""
	return employee, nil
}

func (s *employeeService) GetByID(id string) (*models.User, error) {
	return s.employeeRepo.FindByID(id)
}

func (s *employeeService) List(query *models.PaginationQuery) ([]models.User, int64, error) {
	return s.employeeRepo.List(query.Page, query.PageSize, query.Search)
}
