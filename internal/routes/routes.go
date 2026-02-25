package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/config"
	"github.com/vinodhini/software-api/internal/controllers"
	"github.com/vinodhini/software-api/internal/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	cfg *config.Config,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	projectController *controllers.ProjectController,
	serviceRequestController *controllers.ServiceRequestController,
	messageController *controllers.MessageController,
	clientController *controllers.ClientController,
	serviceTypeController *controllers.ServiceTypeController,
	employeeController *controllers.EmployeeController,
) {
	api := router.Group("/api")

	// Public route for active service types (accessible by clients)
	api.GET("/service-types", serviceTypeController.GetAll)

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// Employee routes
		employees := protected.Group("/employees")
		{
			employees.POST("", middleware.RoleMiddleware("admin"), employeeController.Create)
			employees.GET("", employeeController.List)
			employees.GET("/:id", employeeController.GetByID)
			employees.PUT("/:id", middleware.RoleMiddleware("admin"), userController.Update)
			employees.PATCH("/:id", middleware.RoleMiddleware("admin"), userController.Patch)
			employees.DELETE("/:id", middleware.RoleMiddleware("admin"), userController.Delete)
		}

		// User routes
		users := protected.Group("/users")
		{
			users.GET("", middleware.RoleMiddleware("admin"), userController.List)
			users.GET("/:id", userController.GetByID)
			users.PUT("/:id", userController.Update)
			users.PATCH("/:id", userController.Patch)
			users.DELETE("/:id", middleware.RoleMiddleware("admin"), userController.Delete)
			users.GET("/dashboard/stats", middleware.RoleMiddleware("employee", "admin", "client"), userController.GetDashboardStats)
		}

		// Client routes
		clients := protected.Group("/clients")
		{
			clients.POST("", clientController.Create)
			clients.GET("", clientController.List)
			clients.GET("/:id", clientController.GetByID)
			clients.PUT("/:id", clientController.Update)
			clients.DELETE("/:id", middleware.RoleMiddleware("admin"), clientController.Delete)
		}

		// Project routes
		projects := protected.Group("/projects")
		{
			projects.POST("", middleware.RoleMiddleware("admin"), projectController.Create)
			projects.GET("", projectController.List)
			projects.GET("/:id", projectController.GetByID)
			projects.PUT("/:id", middleware.RoleMiddleware("admin", "employee"), projectController.Update)
			projects.DELETE("/:id", middleware.RoleMiddleware("admin"), projectController.Delete)
			projects.POST("/:id/assign", middleware.RoleMiddleware("admin"), projectController.AssignEmployees)
			projects.PATCH("/:id/progress", middleware.RoleMiddleware("admin", "employee", "client"), projectController.UpdateProgress)
			projects.GET("/:id/messages", messageController.ListByProject)
		}

		// Service request routes
		serviceRequests := protected.Group("/service-requests")
		{
			serviceRequests.POST("", middleware.RoleMiddleware("client"), serviceRequestController.Create)
			serviceRequests.GET("", serviceRequestController.List)
			serviceRequests.GET("/:id", serviceRequestController.GetByID)
			serviceRequests.PUT("/:id", middleware.RoleMiddleware("admin", "employee"), serviceRequestController.Update)
			serviceRequests.DELETE("/:id", middleware.RoleMiddleware("admin"), serviceRequestController.Delete)
			serviceRequests.POST("/:id/approve", middleware.RoleMiddleware("admin"), serviceRequestController.Approve)
			serviceRequests.POST("/:id/reject", middleware.RoleMiddleware("admin"), serviceRequestController.Reject)
		}

		// Service type routes (admin only for management)
		serviceTypes := protected.Group("/service-types")
		{
			serviceTypes.POST("", middleware.RoleMiddleware("admin"), serviceTypeController.Create)
			serviceTypes.GET("/:id", serviceTypeController.GetByID)
			serviceTypes.PUT("/:id", middleware.RoleMiddleware("admin"), serviceTypeController.Update)
			serviceTypes.DELETE("/:id", middleware.RoleMiddleware("admin"), serviceTypeController.Delete)
		}

		// Message routes
		messages := protected.Group("/messages")
		{
			messages.GET("", messageController.List)
			messages.POST("", messageController.Create)
			messages.GET("/:id", messageController.GetByID)
			messages.DELETE("/:id", messageController.Delete)
		}
	}
}
