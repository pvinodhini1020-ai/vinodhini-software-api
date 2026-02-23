package services

import (
	"errors"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type ServiceTypeService struct {
	repository *repositories.ServiceTypeRepository
}

func NewServiceTypeService(repository *repositories.ServiceTypeRepository) *ServiceTypeService {
	return &ServiceTypeService{
		repository: repository,
	}
}

func (s *ServiceTypeService) Create(req *models.CreateServiceTypeRequest) (*models.ServiceType, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	serviceType := &models.ServiceType{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	err := s.repository.Create(serviceType)
	if err != nil {
		return nil, err
	}

	return serviceType, nil
}

func (s *ServiceTypeService) GetByID(id string) (*models.ServiceType, error) {
	return s.repository.GetByID(id)
}

func (s *ServiceTypeService) GetAll(status *string) ([]models.ServiceType, error) {
	return s.repository.GetAll(status)
}

func (s *ServiceTypeService) GetActive() ([]models.ServiceType, error) {
	activeStatus := string(models.StatusActive)
	return s.repository.GetAll(&activeStatus)
}

func (s *ServiceTypeService) Update(id string, req *models.UpdateServiceTypeRequest) (*models.ServiceType, error) {
	existingServiceType, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existingServiceType.Name = req.Name
	}
	if req.Description != "" {
		existingServiceType.Description = req.Description
	}
	if req.Status != "" {
		existingServiceType.Status = req.Status
	}

	err = s.repository.Update(id, existingServiceType)
	if err != nil {
		return nil, err
	}

	return existingServiceType, nil
}

func (s *ServiceTypeService) Delete(id string) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}
