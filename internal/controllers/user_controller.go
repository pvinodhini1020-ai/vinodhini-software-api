package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [get]
func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	user, err := c.userService.GetByID(id, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employees can only view their own profile" || 
		   err.Error() == "access denied: clients can only view their own profile" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User retrieved successfully", user)
}

// @Summary Update user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "Update User Request"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [put]
func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.Update(id, &req, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "access denied: employees can only update their own profile" || 
		   err.Error() == "access denied: clients can only update their own profile" ||
		   err.Error() == "access denied: employees cannot modify role, department, or salary" ||
		   err.Error() == "access denied: clients cannot modify role or company" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User updated successfully", user)
}

// @Summary Update user (partial)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "Update User Request"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [patch]
func (c *UserController) Patch(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")
	
	// Debug logging
	fmt.Printf("PATCH request for user ID: %s\n", id)

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Debug logging
	fmt.Printf("Update request: %+v\n", req)

	user, err := c.userService.Update(id, &req, userID.(string), userRole.(string))
	if err != nil {
		fmt.Printf("Update error: %v\n", err)
		if err.Error() == "access denied: employees can only update their own profile" || 
		   err.Error() == "access denied: clients can only update their own profile" ||
		   err.Error() == "access denied: employees cannot modify role, department, or salary" ||
		   err.Error() == "access denied: clients cannot modify role or company" {
			utils.ErrorResponse(ctx, http.StatusForbidden, err.Error())
		} else {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update user")
		}
		return
	}
	
	// Debug logging
	fmt.Printf("Updated user: %+v\n", user)

	utils.SuccessResponse(ctx, http.StatusOK, "User updated successfully", user)
}

// @Summary Delete user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userService.Delete(id); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User deleted successfully", nil)
}

// @Summary List users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search term"
// @Param role query string false "Filter by role"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/users [get]
func (c *UserController) List(ctx *gin.Context) {
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

	role := ctx.Query("role")
	users, total, err := c.userService.List(&query, role)
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

	utils.PaginatedSuccessResponse(ctx, http.StatusOK, users, pagination)
}

// @Summary Get dashboard statistics
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/users/dashboard/stats [get]
func (c *UserController) GetDashboardStats(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	userRole, _ := ctx.Get("user_role")

	stats, err := c.userService.GetDashboardStats(userID.(string), userRole.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Dashboard statistics retrieved successfully", stats)
}
