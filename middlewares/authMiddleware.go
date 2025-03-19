package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitconnect-backend/utils"
)

// AuthMiddleware verifies the JWT token in the request header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Ensure token is in "Bearer <token>" format
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Pass user ID to the request context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

