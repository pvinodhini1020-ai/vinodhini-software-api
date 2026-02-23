package services

import (
	"errors"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type ServiceRequestService interface {
	Create(clientID string, req *models.CreateServiceRequestRequest) (*models.ServiceRequest, error)
	GetByID(id string) (*models.ServiceRequest, error)
	Update(id string, req *models.UpdateServiceRequestRequest) (*models.ServiceRequest, error)
	Delete(id string) error
	List(query *models.PaginationQuery, clientID *string) ([]models.ServiceRequest, int64, error)
}

type serviceRequestService struct {
	serviceRequestRepo repositories.ServiceRequestRepository
}

func NewServiceRequestService(serviceRequestRepo repositories.ServiceRequestRepository) ServiceRequestService {
	return &serviceRequestService{serviceRequestRepo: serviceRequestRepo}
}

func (s *serviceRequestService) Create(clientID string, req *models.CreateServiceRequestRequest) (*models.ServiceRequest, error) {
	serviceRequest := &models.ServiceRequest{
		Title:       req.Title,
		Description: req.Description,
		ClientID:    clientID,
		Status:      models.StatusPending,
	}

	if req.ProjectID != nil {
		serviceRequest.ProjectID = req.ProjectID
	}

	if err := s.serviceRequestRepo.Create(serviceRequest); err != nil {
		return nil, err
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
