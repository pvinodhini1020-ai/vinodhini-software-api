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
	message, err := c.messageService.Create(userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Message created successfully", message)
}

func (c *MessageController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	message, err := c.messageService.GetByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Message retrieved successfully", message)
}

func (c *MessageController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.messageService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Message deleted successfully", nil)
}

func (c *MessageController) ListByProject(ctx *gin.Context) {
	projectID := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	messages, total, err := c.messageService.ListByProject(projectID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
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
