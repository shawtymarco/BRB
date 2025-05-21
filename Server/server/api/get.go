package api

import "github.com/gin-gonic/gin"

func initGetRequests(router *gin.Engine) {
	router.GET("/hello", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Gin!",
		})
	})
}
