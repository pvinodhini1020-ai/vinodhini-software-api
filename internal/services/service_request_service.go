package services

import (
	"errors"
	"fmt"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type ServiceRequestService interface {
	Create(clientID string, req *models.CreateServiceRequestRequest) (*models.ServiceRequest, error)
	GetByID(id string) (*models.ServiceRequest, error)
	Update(id string, req *models.UpdateServiceRequestRequest) (*models.ServiceRequest, error)
	Delete(id string) error
	List(query *models.PaginationQuery, clientID *string) ([]models.ServiceRequest, int64, error)
	Approve(id string, employeeIDs *[]string) (*models.Project, error)
	Reject(id string) error
}

type serviceRequestService struct {
	serviceRequestRepo repositories.ServiceRequestRepository
	projectRepo         repositories.ProjectRepository
	counterRepo         repositories.CounterRepository
}

func NewServiceRequestService(serviceRequestRepo repositories.ServiceRequestRepository, projectRepo repositories.ProjectRepository, counterRepo repositories.CounterRepository) ServiceRequestService {
	return &serviceRequestService{
		serviceRequestRepo: serviceRequestRepo,
		projectRepo:         projectRepo,
		counterRepo:         counterRepo,
	}
}

func (s *serviceRequestService) Create(clientID string, req *models.CreateServiceRequestRequest) (*models.ServiceRequest, error) {
	// Validate request
	if req.Title == "" {
		return nil, fmt.Errorf("service request title is required")
	}
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	// Generate next service ID sequence
	sequence, err := s.counterRepo.GetNextSequence("service_request_counter")
	if err != nil {
		return nil, fmt.Errorf("failed to generate service ID: %w", err)
	}

	serviceID := fmt.Sprintf("SERVICE%02d", sequence)

	serviceRequest := &models.ServiceRequest{
		ID:          serviceID,
		Title:       req.Title,
		Description: req.Description,
		ClientID:    clientID,
		Status:      models.StatusPending,
	}

	if req.ProjectID != nil {
		serviceRequest.ProjectID = req.ProjectID
	}

	if err := s.serviceRequestRepo.Create(serviceRequest); err != nil {
		return nil, fmt.Errorf("failed to create service request: %w", err)
	}

	return s.serviceRequestRepo.FindByID(serviceRequest.ID)
}

func (s *serviceRequestService) GetByID(id string) (*models.ServiceRequest, error) {
	return s.serviceRequestRepo.FindByID(id)
}

func (s *serviceRequestService) Update(id string, req *models.UpdateServiceRequestRequest) (*models.ServiceRequest, error) {
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	if req.Title != "" {
		serviceRequest.Title = req.Title
	}
	if req.Description != "" {
		serviceRequest.Description = req.Description
	}
	if req.Status != "" {
		serviceRequest.Status = req.Status
	}
	if req.ProjectID != nil {
		serviceRequest.ProjectID = req.ProjectID
	}

	if err := s.serviceRequestRepo.Update(serviceRequest); err != nil {
		return nil, err
	}

	return serviceRequest, nil
}

func (s *serviceRequestService) Delete(id string) error {
	_, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return errors.New("service request not found")
	}

	return s.serviceRequestRepo.Delete(id)
}

func (s *serviceRequestService) List(query *models.PaginationQuery, clientID *string) ([]models.ServiceRequest, int64, error) {
	return s.serviceRequestRepo.List(query.Page, query.PageSize, query.Search, string(query.Status), clientID)
}

func (s *serviceRequestService) Approve(id string, employeeIDs *[]string) (*models.Project, error) {
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("service request not found")
	}

	if serviceRequest.Status != models.StatusPending {
		return nil, errors.New("service request is not pending")
	}

	// Generate next project ID sequence for the new project
	sequence, err := s.counterRepo.GetNextSequence("project_counter")
	if err != nil {
		return nil, fmt.Errorf("failed to generate project ID: %w", err)
	}

	projectID := fmt.Sprintf("PROJECT%02d", sequence)

	// Create project from service request
	project := &models.Project{
		ID:          projectID,
		Name:        serviceRequest.Title,
		Description: serviceRequest.Description,
		ClientID:    serviceRequest.ClientID,
		Status:      models.StatusActive,
		EmployeeIDs: *employeeIDs,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Update service request status and link to project
	serviceRequest.Status = models.StatusActive
	serviceRequest.ProjectID = &project.ID
	if err := s.serviceRequestRepo.Update(serviceRequest); err != nil {
		return nil, fmt.Errorf("failed to update service request: %w", err)
	}

	return s.projectRepo.FindByID(project.ID)
}

func (s *serviceRequestService) Reject(id string) error {
	serviceRequest, err := s.serviceRequestRepo.FindByID(id)
	if err != nil {
		return errors.New("service request not found")
	}

	if serviceRequest.Status != models.StatusPending {
		return errors.New("service request is not pending")
	}

	serviceRequest.Status = models.StatusRejected
	return s.serviceRequestRepo.Update(serviceRequest)
}
