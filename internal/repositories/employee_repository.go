package repositories

import (
	"github.com/vinodhini/software-api/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository interface {
	Create(employee *models.User) error
	FindByID(id string) (*models.User, error)
	List(page, pageSize int, search string) ([]models.User, int64, error)
}

type employeeRepository struct {
	userRepo UserRepository
}

func NewEmployeeRepository(db *mongo.Database) EmployeeRepository {
	return &employeeRepository{
		userRepo: NewUserRepository(db),
	}
}

func (r *employeeRepository) Create(employee *models.User) error {
	employee.Role = models.RoleEmployee
	return r.userRepo.Create(employee)
}

func (r *employeeRepository) FindByID(id string) (*models.User, error) {
	return r.userRepo.FindByID(id)
}

func (r *employeeRepository) List(page, pageSize int, search string) ([]models.User, int64, error) {
	return r.userRepo.List(page, pageSize, search, string(models.RoleEmployee))
}
