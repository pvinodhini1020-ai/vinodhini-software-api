package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinodhini/software-api/pkg/utils"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Err)

			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		}
	}
}
