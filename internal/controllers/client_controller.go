package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type ClientController struct {
	clientService services.ClientService
}

func NewClientController(clientService services.ClientService) *ClientController {
	return &ClientController{
		clientService: clientService,
	}
}

type ClientResponse struct {
	ID       string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Hide     bool   `json:"hide"`
}

// @Summary Create client
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateClientRequest true "Create Client Request"
// @Success 201 {object} utils.Response
// @Router /api/clients [post]
func (c *ClientController) Create(ctx *gin.Context) {
	var req models.CreateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "client"
	}

	client, err := c.clientService.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create client")
		return
	}

	response := ClientResponse{
		ID:      client.UserID,
		Name:    client.Name,
		Email:   client.Email,
		Phone:   client.Phone,
		Company: client.Company,
		Address: client.Address,
		Role:    string(client.Role),
		Status:  client.Status,
		Hide:    client.Hide,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Client created successfully", response)
}

// @Summary List clients
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search term"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/clients [get]
func (c *ClientController) List(ctx *gin.Context) {
	query := &models.PaginationQuery{}
	if err := ctx.ShouldBindQuery(query); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	clients, total, err := c.clientService.List(query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch clients")
		return
	}

	var response []ClientResponse
	for _, client := range clients {
		response = append(response, ClientResponse{
			ID:      client.UserID,
			Name:    client.Name,
			Email:   client.Email,
			Phone:   client.Phone,
			Company: client.Company,
			Address: client.Address,
			Role:    string(client.Role),
			Status:  client.Status,
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

// @Summary Get client by ID
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} utils.Response
// @Router /api/clients/{id} [get]
func (c *ClientController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Client ID is required")
		return
	}

	client, err := c.clientService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Client not found")
		return
	}

	response := ClientResponse{
		ID:      client.UserID,
		Name:    client.Name,
		Email:   client.Email,
		Phone:   client.Phone,
		Company: client.Company,
		Address: client.Address,
		Role:    string(client.Role),
		Status:  client.Status,
		Hide:    client.Hide,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Client retrieved successfully", response)
}

// @Summary Update client
// @Tags clients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body models.UpdateUserRequest true "Update Client Request"
// @Success 200 {object} utils.Response
// @Router /api/clients/{id} [put]
func (c *ClientController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Printf("Update request for client ID: %s\n", id)

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	fmt.Printf("Update request payload: %+v\n", req)

	client, err := c.clientService.Update(id, &req)
	if err != nil {
		fmt.Printf("Service update error: %v\n", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update client: "+err.Error())
		return
	}

	response := ClientResponse{
		ID:      client.UserID,
		Name:    client.Name,
		Email:   client.Email,
		Phone:   client.Phone,
		Company: client.Company,
		Address: client.Address,
		Role:    string(client.Role),
		Status:  client.Status,
		Hide:    client.Hide,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Client updated successfully", response)
}

// @Summary Delete client
// @Tags clients
// @Security BearerAuth
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} utils.Response
// @Router /api/clients/{id} [delete]
func (c *ClientController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Client ID is required")
		return
	}

	if err := c.clientService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete client")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Client deleted successfully", nil)
}
