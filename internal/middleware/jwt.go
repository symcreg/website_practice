package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"test/internal/utility"
	"time"
)

func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		token := auth[len("Bearer"):]
		token = strings.TrimSpace(token)
		claims, err := utility.ParseToken(token)
		if err != nil {
			//c.JSON(401, gin.H{"error": "Invalid token"})
			//c.Abort()
			c.Set("email", "")
			c.Next()
			return
		}
		// Check if token is valid
		if !utility.IsTokenValid(claims, token) {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Update token if it is about to expire in 5 minutes
		if time.Now().Add(time.Minute*5).Unix() > claims.ExpiresAt.Unix() {
			token, err := utility.GenerateToken(claims.Email)
			if err != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
			c.Header("Authorization", "Bearer "+token) // Update token in response header
		}
		c.Set("email", claims.Email)
		c.Next()
	}
}
