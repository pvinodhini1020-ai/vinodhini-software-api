package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/internal/models"
	"github.com/vinodhini/software-api/internal/services"
	"github.com/vinodhini/software-api/pkg/utils"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} utils.Response
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User registered successfully", user)
}

// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} utils.Response
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}
