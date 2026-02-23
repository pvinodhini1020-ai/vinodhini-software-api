package services

import (
	"errors"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type ProjectService interface {
	Create(req *models.CreateProjectRequest) (*models.Project, error)
	GetByID(id string) (*models.Project, error)
	Update(id string, req *models.UpdateProjectRequest) (*models.Project, error)
	Delete(id string) error
	List(query *models.PaginationQuery, clientID *string) ([]models.Project, int64, error)
	AssignEmployees(projectID string, req *models.AssignEmployeesRequest) error
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) Create(req *models.CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		ClientID:    req.ClientID,
		Status:      models.StatusPending,
	}

	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) GetByID(id string) (*models.Project, error) {
	return s.projectRepo.FindByID(id)
}

func (s *projectService) Update(id string, req *models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("project not found")
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

func (s *projectService) List(query *models.PaginationQuery, clientID *string) ([]models.Project, int64, error) {
	return s.projectRepo.List(query.Page, query.PageSize, query.Search, string(query.Status), clientID)
}

func (s *projectService) AssignEmployees(projectID string, req *models.AssignEmployeesRequest) error {
	return s.projectRepo.AssignEmployees(projectID, req.EmployeeIDs)
}
