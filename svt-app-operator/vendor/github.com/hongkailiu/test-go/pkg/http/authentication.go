package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hongkailiu/test-go/pkg/swagger/swagger/models"
)

// AuthenticationMiddleware handles if the context is not authenticated
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthenticated(c) {
			msg := "authentication is required"
			c.JSON(http.StatusUnauthorized, models.Error{Code: int64(http.StatusUnauthorized), Message: &msg})
			c.Abort()
			return
		}
	}
}

// IsAuthenticated returns true if the context is authenticated
func IsAuthenticated(c *gin.Context) bool {
	return getKeyInSession(c, "username") != nil
}
