package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     Role   `json:"role" binding:"required,oneof=admin employee client"`
}

type CreateEmployeeRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Role     string `json:"role"`
	Department string `json:"department" binding:"required"`
	Salary   int    `json:"salary" binding:"required,min=0"`
	Password string `json:"password" binding:"required,min=6"`
	Status   string `json:"status" binding:"required,oneof=active inactive"`
}

type CreateClientRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Company  string `json:"company" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Role     string `json:"role"`
	Password string `json:"password" binding:"required,min=6"`
	Status   string `json:"status" binding:"required,oneof=active inactive"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateUserRequest struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty" binding:"omitempty,email"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty" binding:"omitempty,oneof=admin employee client"`
	Department string `json:"department,omitempty"`
	Company   string `json:"company,omitempty"`
	Address   string `json:"address,omitempty"`
	Salary    int    `json:"salary,omitempty"`
	Password  string `json:"password,omitempty"`
	Status    string `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
	Hide      *bool  `json:"hide,omitempty"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ClientID    string `json:"client_id" binding:"required"`
	Status      Status `json:"status" binding:"omitempty,oneof=active pending completed in_progress"`
	EmployeeIDs []string `json:"employee_ids,omitempty"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      Status `json:"status,omitempty" binding:"omitempty,oneof=active pending completed rejected in_progress"`
}

type AssignEmployeesRequest struct {
	EmployeeIDs []string `json:"employee_ids" binding:"required"`
}

type UpdateProjectProgressRequest struct {
	Progress int `json:"progress" binding:"required,min=0,max=100"`
}

type CreateServiceRequestRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	ProjectID   *string `json:"project_id,omitempty"`
}

type UpdateServiceRequestRequest struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Status      Status  `json:"status,omitempty" binding:"omitempty,oneof=active pending completed rejected"`
	ProjectID   *string `json:"project_id,omitempty"`
}

type CreateMessageRequest struct {
	Content   string `json:"content" binding:"required"`
	ProjectID string `json:"project_id" binding:"required"`
}

type PaginationQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"`
	Status   string `form:"status" binding:"omitempty,oneof=active pending completed rejected"`
}

type CreateServiceTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Status      Status `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateServiceTypeRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      Status `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
}
