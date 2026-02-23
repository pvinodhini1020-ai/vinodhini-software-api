package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type ServiceRequestController struct {
	serviceRequestService services.ServiceRequestService
}

func NewServiceRequestController(serviceRequestService services.ServiceRequestService) *ServiceRequestController {
	return &ServiceRequestController{serviceRequestService: serviceRequestService}
}

func (c *ServiceRequestController) Create(ctx *gin.Context) {
	var req models.CreateServiceRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	serviceRequest, err := c.serviceRequestService.Create(userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Service request created successfully", serviceRequest)
}

func (c *ServiceRequestController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	serviceRequest, err := c.serviceRequestService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service request retrieved successfully", serviceRequest)
}

func (c *ServiceRequestController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateServiceRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	serviceRequest, err := c.serviceRequestService.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service request updated successfully", serviceRequest)
}

func (c *ServiceRequestController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.serviceRequestService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service request deleted successfully", nil)
}

func (c *ServiceRequestController) List(ctx *gin.Context) {
	var query models.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	var clientID *string
	userRole, _ := ctx.Get("user_role")
	if userRole == "client" {
		userIDVal, _ := ctx.Get("user_id")
		uid := userIDVal.(string)
		clientID = &uid
	}

	requests, total, err := c.serviceRequestService.List(&query, clientID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int(total) / query.PageSize
	if int(total)%query.PageSize != 0 {
		totalPages++
	}

	pagination := utils.Pagination{
		Page:      query.Page,
		PageSize:  query.PageSize,
		Total:     total,
		TotalPage: totalPages,
	}

	utils.PaginatedSuccessResponse(ctx, http.StatusOK, requests, pagination)
}

func (c *ServiceRequestController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		EmployeeIDs []string `json:"employee_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	project, err := c.serviceRequestService.Approve(id, &req.EmployeeIDs)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service request approved and project created successfully", project)
}

func (c *ServiceRequestController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.serviceRequestService.Reject(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service request rejected successfully", nil)
}
