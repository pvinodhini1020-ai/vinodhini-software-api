package services

import (
	"errors"
	"testing"
	"time"

	"github.com/vinodhini/software-api/config"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/pkg/utils"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.UserID] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) FindByUserID(userID string) (*models.User, error) {
	return m.FindByID(userID)
}

func (m *MockUserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	m.users[user.UserID] = user
	return nil
}

func (m *MockUserRepository) Delete(id string) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(page, pageSize int, search string, role string) ([]models.User, int64, error) {
	var users []models.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, int64(len(users)), nil
}

func (m *MockUserRepository) GetNextUserID() (string, error) {
	return "USER01", nil
}

func TestLogin_ActiveUser(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	jwtExpiry, _ := time.ParseDuration("24h")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: jwtExpiry,
		},
	}
	
	authService := NewAuthService(mockRepo, cfg)
	
	// Create a test user with active status
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := &models.User{
		UserID:   "USER01",
		Email:    "test@example.com",
		Password: hashedPassword,
		Name:     "Test User",
		Role:     models.RoleClient,
		Status:   "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(testUser)
	
	// Test login with active user
	loginReq := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	response, err := authService.Login(loginReq)
	
	// Assertions
	if err != nil {
		t.Errorf("Expected no error for active user, got: %v", err)
	}
	
	if response == nil {
		t.Error("Expected response, got nil")
	}
	
	if response.Token == "" {
		t.Error("Expected token in response")
	}
	
	if response.User.Email != testUser.Email {
		t.Errorf("Expected user email %s, got %s", testUser.Email, response.User.Email)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	jwtExpiry, _ := time.ParseDuration("24h")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: jwtExpiry,
		},
	}
	
	authService := NewAuthService(mockRepo, cfg)
	
	// Create a test user with inactive status
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := &models.User{
		UserID:   "USER02",
		Email:    "inactive@example.com",
		Password: hashedPassword,
		Name:     "Inactive User",
		Role:     models.RoleClient,
		Status:   "inactive",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(testUser)
	
	// Test login with inactive user
	loginReq := &models.LoginRequest{
		Email:    "inactive@example.com",
		Password: "password123",
	}
	
	response, err := authService.Login(loginReq)
	
	// Assertions
	if err == nil {
		t.Error("Expected error for inactive user, got nil")
	}
	
	if response != nil {
		t.Error("Expected no response for inactive user, got response")
	}
	
	expectedError := "account is inactive. Please contact your system administrator to activate your account"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLogin_UserWithoutStatus(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	jwtExpiry, _ := time.ParseDuration("24h")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: jwtExpiry,
		},
	}
	
	authService := NewAuthService(mockRepo, cfg)
	
	// Create a test user without status (should be treated as active)
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := &models.User{
		UserID:   "USER03",
		Email:    "nostatus@example.com",
		Password: hashedPassword,
		Name:     "No Status User",
		Role:     models.RoleClient,
		// Status is not set (empty string)
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.Create(testUser)
	
	// Test login with user without status
	loginReq := &models.LoginRequest{
		Email:    "nostatus@example.com",
		Password: "password123",
	}
	
	response, err := authService.Login(loginReq)
	
	// Assertions
	if err != nil {
		t.Errorf("Expected no error for user without status, got: %v", err)
	}
	
	if response == nil {
		t.Error("Expected response, got nil")
	}
	
	if response.Token == "" {
		t.Error("Expected token in response")
	}
}

func TestRegister_DefaultActiveStatus(t *testing.T) {
	// Setup
	mockRepo := NewMockUserRepository()
	jwtExpiry, _ := time.ParseDuration("24h")
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expiry: jwtExpiry,
		},
	}
	
	authService := NewAuthService(mockRepo, cfg)
	
	// Test user registration
	registerReq := &models.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
		Role:     models.RoleClient,
	}
	
	user, err := authService.Register(registerReq)
	
	// Assertions
	if err != nil {
		t.Errorf("Expected no error during registration, got: %v", err)
	}
	
	if user == nil {
		t.Error("Expected user, got nil")
	}
	
	if user.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", user.Status)
	}
}
