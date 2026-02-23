package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type ServiceTypeController struct {
	serviceTypeService *services.ServiceTypeService
}

func NewServiceTypeController(serviceTypeService *services.ServiceTypeService) *ServiceTypeController {
	return &ServiceTypeController{serviceTypeService: serviceTypeService}
}

func (c *ServiceTypeController) Create(ctx *gin.Context) {
	var req models.CreateServiceTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	serviceType, err := c.serviceTypeService.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Service type created successfully", serviceType)
}

func (c *ServiceTypeController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	serviceType, err := c.serviceTypeService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service type retrieved successfully", serviceType)
}

func (c *ServiceTypeController) GetAll(ctx *gin.Context) {
	status := ctx.Query("status")

	var serviceTypes []models.ServiceType
	var err error

	if status == "active" {
		serviceTypes, err = c.serviceTypeService.GetActive()
	} else {
		var statusPtr *string
		if status != "" {
			statusPtr = &status
		}
		serviceTypes, err = c.serviceTypeService.GetAll(statusPtr)
	}

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service types retrieved successfully", serviceTypes)
}

func (c *ServiceTypeController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.UpdateServiceTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	serviceType, err := c.serviceTypeService.Update(id, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service type updated successfully", serviceType)
}

func (c *ServiceTypeController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.serviceTypeService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Service type deleted successfully", nil)
}
