package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

func GlobalCommands() {
	cmd.Register(cmd.New("ping", "To check a specific player's ping", nil, PingCommand{}))
	cmd.Register(cmd.New("color", "To show list of colors", nil, ColorCommand{}))
	cmd.Register(cmd.New("debug", "To debug certain code", nil, DebugCommand{}))
	cmd.Register(cmd.New("gamemode", "To set a specific player's game mode", nil, GamemodeCommand{}))
	cmd.Register(cmd.New("rank", "To set a specific player's rank", nil, RankCommand{}))
}

func GameCommands() {
	cmd.Register(cmd.New("hub", "To teleport back to hub", nil, HubCommand{}))
	cmd.Register(cmd.New("game", "To manage the game you are currently in", nil, GameStartCommand{}))
}
