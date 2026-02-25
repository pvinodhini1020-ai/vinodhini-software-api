package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type EmployeeController struct {
	employeeService services.EmployeeService
}

func NewEmployeeController(employeeService services.EmployeeService) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
	}
}

type EmployeeResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Salary     int    `json:"salary,omitempty"`
}

func (c *EmployeeController) Create(ctx *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "employee"
	}

	employee, err := c.employeeService.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create employee")
		return
	}

	response := EmployeeResponse{
		ID:         employee.UserID,
		Name:       employee.Name,
		Email:      employee.Email,
		Role:       string(employee.Role),
		Status:     employee.Status,
		Phone:      employee.Phone,
		Department:  employee.Department,
		Salary:      employee.Salary,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Employee created successfully", response)
}

func (c *EmployeeController) List(ctx *gin.Context) {
	query := &models.PaginationQuery{}
	if err := ctx.ShouldBindQuery(query); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	employees, total, err := c.employeeService.List(query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}

	var response []EmployeeResponse
	for _, emp := range employees {
		response = append(response, EmployeeResponse{
			ID:         emp.UserID,
			Name:       emp.Name,
			Email:      emp.Email,
			Role:       string(emp.Role),
			Status:     emp.Status,
			Phone:      emp.Phone,
			Department:  emp.Department,
			Salary:      emp.Salary,
		})
	}

	pagination := utils.Pagination{
		Page:      query.Page,
		PageSize:  query.PageSize,
		Total:     total,
		TotalPage: int((total + int64(query.PageSize) - 1) / int64(query.PageSize)),
	}
	utils.PaginatedSuccessResponse(ctx, http.StatusOK, response, pagination)
}

func (c *EmployeeController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Employee ID is required")
		return
	}

	employee, err := c.employeeService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Employee not found")
		return
	}

	response := EmployeeResponse{
		ID:         employee.UserID,
		Name:       employee.Name,
		Email:      employee.Email,
		Role:       string(employee.Role),
		Status:     employee.Status,
		Phone:      employee.Phone,
		Department:  employee.Department,
		Salary:      employee.Salary,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Employee retrieved successfully", response)
}
