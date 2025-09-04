package api

import (
	"net/http"
	"server/server"
	"server/server/database"
	"server/server/game"
	"server/server/games/bedwars"
	"strconv"
	"time"

	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
)

func initGetRequests(rg *gin.RouterGroup) {
	now := time.Now()
	rg.GET("/connect", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API successfully connected",
			"time":    now.Second(),
		})
	})

	rg.GET("/players", jwtAuthMiddleware(), func(c *gin.Context) {
		players, err := server.Database.FindAllPlayers()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": players,
		})
	})

	rg.GET("/players/:player", jwtAuthMiddleware(), func(c *gin.Context) {
		username := c.Param("player")
		pd, err := server.Database.FindPlayerByName(username, &database.PlayerNameSearchOpts{CaseInsensitive: false, PartialMatch: false})
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
		isCustom := lo.If(c.Query("custom") == "1", true).Else(false)
		desiredMap := c.Query("map")

		g := bedwars.NewBedWars(game.TypeBedWars, teamSize, teamCount, isCustom, desiredMap)

		c.JSON(http.StatusOK, gin.H{
			"id": g.ID().String(),
		})
	})

	rg.GET("/games/pending-termination", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ids": bedwars.GamesToTerminate,
		})
	})
}
