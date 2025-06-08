package api

import (
	"net/http"
	"server/server"
	"server/server/database"
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/utils"
	"strconv"

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

	rg.GET("/registered-players/:userid", jwtAuthMiddleware(), func(c *gin.Context) {
		id := c.Param("userid")
		pd, err := server.Database.FindPlayerByDiscordID(id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"registered": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"registered": true, "isTouch": pd.IsTouch(), "data": pd})
	})

	rg.GET("/games/create", jwtAuthMiddleware(), func(c *gin.Context) {
		teamSize, _ := strconv.Atoi(c.Query("teamSize"))
		teamCount, _ := strconv.Atoi(c.Query("teamCount"))
		isCustom := utils.Question(c.Query("custom") == "1", true, false)

		g := bedwars.NewBedWars(game.TypeBedWars, teamSize, teamCount, isCustom)

		c.JSON(http.StatusOK, gin.H{
			"id": g.ID().String(),
		})
	})
}
