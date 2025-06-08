package api

import (
	"fmt"
	"net/http"
	"server/server"
	"server/server/games/bedwars"
	"server/server/utils"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var RegistrationCodes = make(map[uuid.UUID]string)

func initPostRequests(rg *gin.RouterGroup) {
	rg.POST("/secure", jwtAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Minecraft API"})
	})

	rg.POST("/verify", jwtAuthMiddleware(), func(c *gin.Context) {
		var body struct {
			Code   string `json:"code"`
			UserID string `json:"userId"`
		}
		utils.Panic(c.BindJSON(&body))

		var playerId uuid.UUID
		var found bool
		for id, code := range RegistrationCodes {
			if body.Code == code {
				playerId = id
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusOK, gin.H{
				"success": 0,
				"message": "The code is either invalid, used by a player or has already expired! You can get a new code by doing /link in the lobby server (eliagic.club).",
			})
			return
		}

		delete(RegistrationCodes, playerId)
		pd := utils.Panics(server.Database.FindPlayer(playerId))
		pd.UserId = body.UserID
		utils.Panic(server.Database.SavePlayer(pd))
		c.JSON(http.StatusOK, gin.H{
			"success":  1,
			"message":  fmt.Sprintf("You successfully linked your Minecraft account `%v` to your Discord account <@%v>!", pd.Username, body.UserID),
			"username": pd.Username,
			"elo":      pd.Statistics.ELO,
		})
	})

	rg.POST("/games/connect-users", jwtAuthMiddleware(), func(c *gin.Context) {
		var body struct {
			Users []string `json:"users"`
			Code  string   `json:"code"`
		}
		utils.Panic(c.BindJSON(&body))
		var bwGame *bedwars.BedWars
		for _, g := range bedwars.Games {
			if g.ID().String() == body.Code {
				bwGame = g
				break
			}
		}
		bwGame.UsersToJoin = body.Users

		c.JSON(http.StatusOK, gin.H{})
	})
}
