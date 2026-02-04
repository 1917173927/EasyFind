package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT 中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// TODO: 验证 Token 逻辑
		// claims, err := utils.ParseToken(parts[1]) ...

		c.Next()
	}
}
