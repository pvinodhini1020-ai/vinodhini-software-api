package services

import (
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type MessageService interface {
	Create(senderID string, req *models.CreateMessageRequest) (*models.Message, error)
	GetByID(id string) (*models.Message, error)
	Delete(id string) error
	ListByProject(projectID string, page, pageSize int) ([]models.Message, int64, error)
}

type messageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) Create(senderID string, req *models.CreateMessageRequest) (*models.Message, error) {
	message := &models.Message{
		Content:   req.Content,
		SenderID:  senderID,
		ProjectID: req.ProjectID,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return s.messageRepo.FindByID(message.ID)
}

func (s *messageService) GetByID(id string) (*models.Message, error) {
	return s.messageRepo.FindByID(id)
}

func (s *messageService) Delete(id string) error {
	return s.messageRepo.Delete(id)
}

func (s *messageService) ListByProject(projectID string, page, pageSize int) ([]models.Message, int64, error) {
	return s.messageRepo.ListByProject(projectID, page, pageSize)
}
