package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// r.Use(middleware.Cors())

	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) { c.JSON(200, gin.H{"message": "login"}) })
			// auth.POST("/register", ...)
		}

		// Items routes
		// items := v1.Group("/items")
		// ...
	}

	return r
}
