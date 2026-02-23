package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(projectService services.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (c *ProjectController) Create(ctx *gin.Context) {
	var req models.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	project, err := c.projectService.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Project created successfully", project)
}

func (c *ProjectController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	project, err := c.projectService.GetByID(id, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only view their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Project retrieved successfully", project)
}

func (c *ProjectController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	var req models.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	project, err := c.projectService.Update(id, &req, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only update their own projects" ||
		   err.Error() == "access denied: employees can only update project status" ||
		   err.Error() == "access denied: clients cannot update project status" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Project updated successfully", project)
}

func (c *ProjectController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.projectService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Project deleted successfully", nil)
}

func (c *ProjectController) List(ctx *gin.Context) {
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

	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")
	
	var clientID *string
	if userRole == "client" {
		uid := userID.(string)
		clientID = &uid
	}

	projects, total, err := c.projectService.List(&query, clientID, userID.(string), userRole.(string))
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

	utils.PaginatedSuccessResponse(ctx, http.StatusOK, projects, pagination)
}

func (c *ProjectController) AssignEmployees(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	var req models.AssignEmployeesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.projectService.AssignEmployees(id, &req, userID.(string), userRole.(string)); err != nil {
		if err.Error() == "access denied: only admins can assign employees to projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Employees assigned successfully", nil)
}

func (c *ProjectController) UpdateProgress(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	var req models.UpdateProjectProgressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	project, err := c.projectService.UpdateProjectProgress(id, &req, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employee not assigned to this project" || 
		   err.Error() == "access denied: client can only update their own projects" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Project progress updated successfully", project)
}
