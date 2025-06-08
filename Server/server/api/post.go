package api

import (
	"fmt"
	"net/http"
	"server/server"
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
}
