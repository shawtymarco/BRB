package api

import (
	"net/http"
	"server/server"
	"server/server/database"

	"github.com/gin-gonic/gin"
)

func initGetRequests(rg *gin.RouterGroup) {
	rg.GET("/connect", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API successfully connected",
		})
	})

	rg.GET("/players/:player", jwtAuthMiddleware(), func(c *gin.Context) {
		username := c.Param("player")
		pd, err := server.Database.FindPlayerFromName(username, &database.PlayerNameSearchOpts{CaseInsensitive: false, PartialMatch: false})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": pd,
		})
	})
}
