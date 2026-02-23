package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/config"
	"github.com/vinodhini/software-api/internal/middleware"
)

func main() {
	cfg := config.Load()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORS.Origins
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Server is running. Connect PostgreSQL to enable full functionality.",
		})
	})

	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Println("Note: PostgreSQL connection required for full API functionality")
	log.Println("Run 'docker-compose up -d' to start PostgreSQL and the full application")
	
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
