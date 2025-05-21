package api

import "github.com/gin-gonic/gin"

func initPostRequests(router *gin.Engine) {
	router.POST("/secure-endpoint", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Minecraft API"})
	})
}
