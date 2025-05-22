package ws

import (
	"encoding/json"
	"log"
	"server/server/utils"

	"github.com/gorilla/websocket"
)

// SendCommandToDiscordBot is the way to send messages to the Discord bot to perform actions on demand.
func SendCommandToDiscordBot(cmd DiscordCommand, onCommandExecute func(response DiscordResponse)) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	utils.Panic(conn.WriteMessage(websocket.TextMessage, utils.Panics(json.Marshal(cmd))))

	go func() {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Read error:", err)
		}

		var response DiscordResponse
		utils.Panic(json.Unmarshal(message, &response))

		onCommandExecute(response)
	}()
}

type DiscordCommand struct {
	Type   string   `json:"type"`
	UserID string   `json:"userId"`
	Args   []string `json:"args"`
}

type DiscordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
