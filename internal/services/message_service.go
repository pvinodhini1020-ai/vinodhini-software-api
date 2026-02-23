package services

import (
	"fmt"

	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/repositories"
)

type MessageService interface {
	Create(senderID string, req *models.CreateMessageRequest, userID string, userRole string) (*models.Message, error)
	GetByID(id string, userID string, userRole string) (*models.Message, error)
	Delete(id string, userID string, userRole string) error
	ListByProject(projectID string, page, pageSize int, userID string, userRole string) ([]models.Message, int64, error)
}

type messageService struct {
	messageRepo  repositories.MessageRepository
	projectRepo  repositories.ProjectRepository
	counterRepo  repositories.CounterRepository
}

func NewMessageService(messageRepo repositories.MessageRepository, counterRepo repositories.CounterRepository, projectRepo repositories.ProjectRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		counterRepo: counterRepo,
		projectRepo: projectRepo,
	}
}

func (s *messageService) Create(senderID string, req *models.CreateMessageRequest, userID string, userRole string) (*models.Message, error) {
	// Validate request
	if req.Content == "" {
		return nil, fmt.Errorf("message content is required")
	}
	if req.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}
	if senderID == "" {
		return nil, fmt.Errorf("sender ID is required")
	}

	// Verify user has access to the project
	project, err := s.projectRepo.FindByID(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Apply role-based access control for messaging
	if userRole == "employee" {
		// Employees can only message projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, fmt.Errorf("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only message their own projects
		if project.ClientID != userID {
			return nil, fmt.Errorf("access denied: client can only message their own projects")
		}
	}

	// Generate next message ID sequence
	sequence, err := s.counterRepo.GetNextSequence("message_counter")
	if err != nil {
		return nil, fmt.Errorf("failed to generate message ID: %w", err)
	}

	messageID := fmt.Sprintf("MESSAGE%02d", sequence)

	message := &models.Message{
		ID:        messageID,
		Content:   req.Content,
		SenderID:  senderID,
		ProjectID: req.ProjectID,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return s.messageRepo.FindByID(message.ID)
}

func (s *messageService) GetByID(id string, userID string, userRole string) (*models.Message, error) {
	message, err := s.messageRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify user has access to the project
	project, err := s.projectRepo.FindByID(message.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only view messages from projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, fmt.Errorf("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only view messages from their own projects
		if project.ClientID != userID {
			return nil, fmt.Errorf("access denied: client can only view their own projects")
		}
	}

	return message, nil
}

func (s *messageService) Delete(id string, userID string, userRole string) error {
	message, err := s.messageRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Verify user has access to the project
	project, err := s.projectRepo.FindByID(message.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found")
	}

	// Apply role-based access control for deletion
	if userRole == "employee" {
		// Employees can only delete their own messages from projects they're assigned to
		if message.SenderID != userID {
			return fmt.Errorf("access denied: can only delete own messages")
		}
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return fmt.Errorf("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only delete their own messages from their own projects
		if message.SenderID != userID {
			return fmt.Errorf("access denied: can only delete own messages")
		}
		if project.ClientID != userID {
			return fmt.Errorf("access denied: client can only delete messages from their own projects")
		}
	}

	return s.messageRepo.Delete(id)
}

func (s *messageService) ListByProject(projectID string, page, pageSize int, userID string, userRole string) ([]models.Message, int64, error) {
	// If projectID is empty, return all messages the user has access to
	if projectID == "" {
		// For now, return empty list as general message listing isn't implemented
		// In a real implementation, you might want to fetch all messages for the user
		return []models.Message{}, 0, nil
	}

	// Verify user has access to the project
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, 0, fmt.Errorf("project not found")
	}

	// Apply role-based access control
	if userRole == "employee" {
		// Employees can only view messages from projects they're assigned to
		isAssigned := false
		for _, empID := range project.EmployeeIDs {
			if empID == userID {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, 0, fmt.Errorf("access denied: employee not assigned to this project")
		}
	} else if userRole == "client" {
		// Clients can only view messages from their own projects
		if project.ClientID != userID {
			return nil, 0, fmt.Errorf("access denied: client can only view their own projects")
		}
	}

	return s.messageRepo.ListByProject(projectID, page, pageSize)
}
