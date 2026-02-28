package middleware

import (
	"easyfind/internal/models"
	"easyfind/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT 中间件：解析 Token 并注入用户信息
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		token := ""
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization header format must be Bearer {token}"})
				c.Abort()
				return
			}
			token = parts[1]
		} else {
			token = strings.TrimSpace(c.Query("token"))
			if token == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization header or token query is required"})
				c.Abort()
				return
			}
		}

		// 验证 Token 逻辑
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid token"})
			c.Abort()
			return
		}

		// 将 claims 存入上下文，供后续 Handler 使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role) // 存入的是 int 类型

		c.Next()
	}
}

// AuthAdmin 要求用户至少是失物招领管理员 (RoleLFAdmin) 或 系统管理员 (RoleSysAdmin)
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Unauthorized"})
			c.Abort()
			return
		}

		r := models.UserRole(role.(int))
		if r != models.RoleLFAdmin && r != models.RoleSysAdmin {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "Forbidden: Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthSysAdmin 要求用户必须是系统管理员 (RoleSysAdmin)
func AuthSysAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Unauthorized"})
			c.Abort()
			return
		}

		r := models.UserRole(role.(int))
		if r != models.RoleSysAdmin {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "Forbidden: System Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
