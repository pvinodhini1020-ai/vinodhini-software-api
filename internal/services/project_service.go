package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type ProjectService interface {
	Create(req *models.CreateProjectRequest) (*models.Project, error)
	GetByID(id string, userID string, userRole string) (*models.Project, error)
	Update(id string, req *models.UpdateProjectRequest, userID string, userRole string) (*models.Project, error)
	Delete(id string) error
	List(query *models.PaginationQuery, clientID *string, userID string, userRole string) ([]models.Project, int64, error)
	AssignEmployees(projectID string, req *models.AssignEmployeesRequest, userID string, userRole string) error
	UpdateProjectProgress(projectID string, req *models.UpdateProjectProgressRequest, userID string, userRole string) (*models.Project, error)
}

type projectService struct {
	projectRepo repositories.ProjectRepository
	counterRepo repositories.CounterRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository, counterRepo repositories.CounterRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		counterRepo: counterRepo,
	}
}

func (s *projectService) Create(req *models.CreateProjectRequest) (*models.Project, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if req.ClientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	// Generate next project ID sequence
	sequence, err := s.counterRepo.GetNextSequence("project_counter")
	if err != nil {
		return nil, fmt.Errorf("failed to generate project ID: %w", err)
	}

	projectID := fmt.Sprintf("PROJECT%02d", sequence)
	
	// Debug: Ensure projectID is not empty
	if projectID == "" {
		return nil, fmt.Errorf("generated project ID is empty")
	}

	project := &models.Project{
		ID:          projectID,
		Name:        req.Name,
		Description: req.Description,
		ClientID:    req.ClientID,
		Status:      models.StatusPending,
		Progress:    0,
		EmployeeIDs: req.EmployeeIDs,
	}

	if req.Status != "" {
		project.Status = req.Status
	}

	// Debug: Log the project before creation
	fmt.Printf("Creating project with ID: %s\n", project.ID)

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Verify the project was created with the correct ID
	createdProject, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created project: %w", err)
	}
	
	// Debug: Log the retrieved project
	fmt.Printf("Retrieved project with ID: %s\n", createdProject.ID)

	return createdProject, nil
}

func (s *projectService) GetByID(id string, userID string, userRole string) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only view projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, errors.New("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only view their own projects
		if project.ClientID != userID {
			return nil, errors.New("access denied: client can only view their own projects")
		}
	}

	return project, nil
}

func (s *projectService) Update(id string, req *models.UpdateProjectRequest, userID string, userRole string) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only update projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, errors.New("access denied: employee not assigned to this project")
		}
		// Employees can only update status, not other fields
		if req.Name != "" || req.Description != "" {
			return nil, errors.New("access denied: employees can only update project status")
		}
	} else if userRole == "client" {
		// Clients can only update their own projects
		if project.ClientID != userID {
			return nil, errors.New("access denied: client can only update their own projects")
		}
		// Clients can only update description, not status or assignment
		if req.Status != "" {
			return nil, errors.New("access denied: clients cannot update project status")
		}
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) Delete(id string) error {
	_, err := s.projectRepo.FindByID(id)
	if err != nil {
		return errors.New("project not found")
	}

	return s.projectRepo.Delete(id)
}

func (s *projectService) List(query *models.PaginationQuery, clientID *string, userID string, userRole string) ([]models.Project, int64, error) {
	// Apply role-based filtering
	if userRole == "employee" {
		// Employees can only see projects they're assigned to
		return s.projectRepo.ListByEmployee(query.Page, query.PageSize, query.Search, string(query.Status), userID)
	} else if userRole == "client" {
		// Clients can only see their own projects
		return s.projectRepo.List(query.Page, query.PageSize, query.Search, string(query.Status), &userID)
	}
	
	// Admins can see all projects
	return s.projectRepo.List(query.Page, query.PageSize, query.Search, string(query.Status), clientID)
}

func (s *projectService) AssignEmployees(projectID string, req *models.AssignEmployeesRequest, userID string, userRole string) error {
	// Only admins can assign employees
	if userRole != "admin" {
		return errors.New("access denied: only admins can assign employees to projects")
	}

	// Get the current project to check existing assignments
	_, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return errors.New("project not found")
	}

	// Validate that all employee IDs exist and are employees
	for _, empID := range req.EmployeeIDs {
		// Additional validation can be added here to verify each employee exists
		// and has the "employee" role
		if empID == "" {
			return errors.New("invalid employee ID provided")
		}
	}

	return s.projectRepo.AssignEmployees(projectID, req.EmployeeIDs)
}

func (s *projectService) UpdateProjectProgress(projectID string, req *models.UpdateProjectProgressRequest, userID string, userRole string) (*models.Project, error) {
	// Get the project to check access
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only update projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, errors.New("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only update their own projects
		if project.ClientID != userID {
			return nil, errors.New("access denied: client can only update their own projects")
		}
	}

	// Validate progress value
	if req.Progress < 0 || req.Progress > 100 {
		return nil, errors.New("progress must be between 0 and 100")
	}

	// Update project progress
	project.Progress = req.Progress
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}
