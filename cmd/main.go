package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/config"
	"github.com/vinodhini/software-api/internal/controllers"
	"github.com/vinodhini/software-api/internal/middleware"
	"github.com/vinodhini/software-api/internal/repositories"
	"github.com/vinodhini/software-api/internal/routes"
	"github.com/vinodhini/software-api/internal/services"
)

// @title Vinodhini Software API
// @version 1.0
// @description Scalable REST API with Clean Architecture
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create indexes
	if err := config.CreateIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	serviceRequestRepo := repositories.NewServiceRequestRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	counterRepo := repositories.NewCounterRepository(db)
	serviceTypeRepo := repositories.NewServiceTypeRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo, projectRepo)
	clientService := services.NewClientService(userRepo)
	projectService := services.NewProjectService(projectRepo, counterRepo)
	serviceRequestService := services.NewServiceRequestService(serviceRequestRepo, projectRepo, counterRepo)
	messageService := services.NewMessageService(messageRepo, counterRepo, projectRepo)
	serviceTypeService := services.NewServiceTypeService(serviceTypeRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	clientController := controllers.NewClientController(clientService)
	projectController := controllers.NewProjectController(projectService)
	serviceRequestController := controllers.NewServiceRequestController(serviceRequestService)
	messageController := controllers.NewMessageController(messageService)
	serviceTypeController := controllers.NewServiceTypeController(serviceTypeService)

	// Setup Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit.Limit, cfg.RateLimit.Window))

	// CORS configuration
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Referer", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup routes
	routes.SetupRoutes(router, cfg, authController, userController, projectController, serviceRequestController, messageController, clientController, serviceTypeController)

	// Server setup
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
