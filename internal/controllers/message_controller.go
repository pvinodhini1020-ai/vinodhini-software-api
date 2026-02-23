package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type MessageController struct {
	messageService services.MessageService
}

func NewMessageController(messageService services.MessageService) *MessageController {
	return &MessageController{messageService: messageService}
}

func (c *MessageController) Create(ctx *gin.Context) {
	var req models.CreateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")
	
	message, err := c.messageService.Create(userID.(string), &req, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only message their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Message created successfully", message)
}

func (c *MessageController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	message, err := c.messageService.GetByID(id, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only view their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Message retrieved successfully", message)
}

func (c *MessageController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	if err := c.messageService.Delete(id, userID.(string), userRole.(string)); err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only delete messages from their own projects" ||
		   err.Error() == "access denied: can only delete own messages" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Message deleted successfully", nil)
}

func (c *MessageController) ListByProject(ctx *gin.Context) {
	projectID := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	messages, total, err := c.messageService.ListByProject(projectID, page, pageSize, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only view their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	pagination := utils.Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPages,
	}

	utils.PaginatedSuccessResponse(ctx, http.StatusOK, messages, pagination)
}

// @Summary List all messages
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/messages [get]
func (c *MessageController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	messages, total, err := c.messageService.ListByProject("", page, pageSize, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only view their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	pagination := utils.Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPages,
	}

	utils.PaginatedSuccessResponse(ctx, http.StatusOK, messages, pagination)
}
