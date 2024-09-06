package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/DomT00T/actico-auth-api/internal/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := auth.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}