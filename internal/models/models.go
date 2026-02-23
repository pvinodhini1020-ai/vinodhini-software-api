package models

import (
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
	RoleClient   Role = "client"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusRejected  Status = "rejected"
	StatusInProgress Status = "in_progress"
)

type User struct {
	UserID    string             `bson:"_id,omitempty" json:"user_id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Name      string             `bson:"name" json:"name"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Role      Role               `bson:"role" json:"role"`
	Department string            `bson:"department,omitempty" json:"department,omitempty"`
	Company   string             `bson:"company,omitempty" json:"company,omitempty"`
	Address   string             `bson:"address,omitempty" json:"address,omitempty"`
	Salary    int                `bson:"salary,omitempty" json:"salary,omitempty"`
	Status    string             `bson:"status,omitempty" json:"status,omitempty"`
	Hide      bool               `bson:"hide,omitempty" json:"hide,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Project struct {
	ID          string   `bson:"_id" json:"id"`
	Name        string   `bson:"name" json:"name"`
	Description string   `bson:"description" json:"description"`
	ClientID    string   `bson:"client_id" json:"client_id"`
	Client      *User    `bson:"-" json:"client,omitempty"`
	Status      Status   `bson:"status" json:"status"`
	Progress    int      `bson:"progress" json:"progress"`
	EmployeeIDs []string `bson:"employee_ids" json:"employee_ids"`
	Employees   []User   `bson:"-" json:"employees,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type ServiceRequest struct {
	ID          string  `bson:"_id" json:"id"`
	Title       string  `bson:"title" json:"title"`
	Description string  `bson:"description" json:"description"`
	ClientID    string  `bson:"client_id" json:"client_id"`
	Client      *User   `bson:"-" json:"client,omitempty"`
	ProjectID   *string `bson:"project_id,omitempty" json:"project_id,omitempty"`
	Project     *Project `bson:"-" json:"project,omitempty"`
	Status      Status  `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type Message struct {
	ID        string    `bson:"_id" json:"id"`
	Content   string    `bson:"content" json:"content"`
	SenderID  string    `bson:"sender_id" json:"sender_id"`
	Sender    *User     `bson:"-" json:"sender,omitempty"`
	ProjectID string    `bson:"project_id" json:"project_id"`
	Project   *Project  `bson:"-" json:"project,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ServiceType struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Status      Status    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
