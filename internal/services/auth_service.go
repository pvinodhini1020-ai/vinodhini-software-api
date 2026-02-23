package services

import (
	"errors"

	"github.com/vinodhini/software-api/config"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
	"github.com/vinodhini/software-api/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate next user ID like USER01, USER02, etc.
	userID, err := s.userRepo.GetNextUserID()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserID:   userID, // This will be used as _id in MongoDB
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Use UserID instead of MongoDB ObjectID for token generation
	token, err := utils.GenerateToken(user.UserID, user.Email, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.Expiry)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
