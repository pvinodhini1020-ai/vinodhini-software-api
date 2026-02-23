package services

import (
	"errors"
	"fmt"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(id string, userID string, userRole string) (*models.User, error)
	Update(id string, req *models.UpdateUserRequest, userID string, userRole string) (*models.User, error)
	Delete(id string) error
	List(query *models.PaginationQuery, role string) ([]models.User, int64, error)
	GetDashboardStats(userID string, userRole string) (map[string]interface{}, error)
}

type userService struct {
	userRepo    repositories.UserRepository
	projectRepo repositories.ProjectRepository
}

func NewUserService(userRepo repositories.UserRepository, projectRepo repositories.ProjectRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

func (s *userService) GetByID(id string, userID string, userRole string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only view their own profile
		if id != userID {
			return nil, errors.New("access denied: employees can only view their own profile")
		}
	} else if userRole == "client" {
		// Clients can only view their own profile
		if id != userID {
			return nil, errors.New("access denied: clients can only view their own profile")
		}
	}

	return user, nil
}

func (s *userService) Update(id string, req *models.UpdateUserRequest, userID string, userRole string) (*models.User, error) {
	fmt.Printf("UserService.Update called with ID: %s, Request: %+v\n", id, req)
	
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		fmt.Printf("User not found: %v\n", err)
		return nil, errors.New("user not found")
	}
	
	fmt.Printf("Found user: %+v\n", user)

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only update their own profile
		if id != userID {
			return nil, errors.New("access denied: employees can only update their own profile")
		}
		// Employees cannot change their role, department, or salary
		if req.Role != "" || req.Department != "" || req.Salary > 0 {
			return nil, errors.New("access denied: employees cannot modify role, department, or salary")
		}
	} else if userRole == "client" {
		// Clients can only update their own profile
		if id != userID {
			return nil, errors.New("access denied: clients can only update their own profile")
		}
		// Clients cannot change their role or company
		if req.Role != "" || req.Company != "" {
			return nil, errors.New("access denied: clients cannot modify role or company")
		}
	}

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
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Company != "" {
		user.Company = req.Company
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

func (s *userService) GetDashboardStats(userID string, userRole string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	if userRole == "employee" {
		// Get employee's assigned projects
		projects, _, err := s.projectRepo.ListByEmployee(1, 1000, "", "", userID)
		if err != nil {
			return nil, err
		}
		
		totalProjects := len(projects)
		activeProjects := 0
		pendingProjects := 0
		completedProjects := 0
		inProgressProjects := 0
		
		for _, project := range projects {
			switch project.Status {
			case models.StatusActive:
				activeProjects++
			case models.StatusPending:
				pendingProjects++
			case models.StatusCompleted:
				completedProjects++
			case models.StatusInProgress:
				inProgressProjects++
			}
		}
		
		stats["assigned_projects"] = totalProjects
		stats["active_projects"] = activeProjects
		stats["pending_projects"] = pendingProjects
		stats["completed_projects"] = completedProjects
		stats["in_progress_projects"] = inProgressProjects
		stats["projects"] = projects
		
	} else if userRole == "client" {
		// Get client's projects
		projects, _, err := s.projectRepo.List(1, 1000, "", "", &userID)
		if err != nil {
			return nil, err
		}
		
		totalProjects := len(projects)
		activeProjects := 0
		pendingProjects := 0
		completedProjects := 0
		inProgressProjects := 0
		
		for _, project := range projects {
			switch project.Status {
			case models.StatusActive:
				activeProjects++
			case models.StatusPending:
				pendingProjects++
			case models.StatusCompleted:
				completedProjects++
			case models.StatusInProgress:
				inProgressProjects++
			}
		}
		
		stats["total_projects"] = totalProjects
		stats["active_projects"] = activeProjects
		stats["pending_projects"] = pendingProjects
		stats["completed_projects"] = completedProjects
		stats["in_progress_projects"] = inProgressProjects
		stats["projects"] = projects
		
	} else if userRole == "admin" {
		// Get all projects for admin overview
		projects, _, err := s.projectRepo.List(1, 1000, "", "", nil)
		if err != nil {
			return nil, err
		}
		
		totalProjects := len(projects)
		activeProjects := 0
		pendingProjects := 0
		completedProjects := 0
		inProgressProjects := 0
		
		for _, project := range projects {
			switch project.Status {
			case models.StatusActive:
				activeProjects++
			case models.StatusPending:
				pendingProjects++
			case models.StatusCompleted:
				completedProjects++
			case models.StatusInProgress:
				inProgressProjects++
			}
		}
		
		// Get user counts
		totalUsers, _, err := s.userRepo.List(1, 1000, "", "")
		if err != nil {
			return nil, err
		}
		
		adminUsers := 0
		employeeUsers := 0
		clientUsers := 0
		
		for _, user := range totalUsers {
			switch user.Role {
			case models.RoleAdmin:
				adminUsers++
			case models.RoleEmployee:
				employeeUsers++
			case models.RoleClient:
				clientUsers++
			}
		}
		
		stats["total_projects"] = totalProjects
		stats["active_projects"] = activeProjects
		stats["pending_projects"] = pendingProjects
		stats["completed_projects"] = completedProjects
		stats["in_progress_projects"] = inProgressProjects
		stats["total_users"] = len(totalUsers)
		stats["admin_users"] = adminUsers
		stats["employee_users"] = employeeUsers
		stats["client_users"] = clientUsers
		stats["projects"] = projects
	}
	
	return stats, nil
}
