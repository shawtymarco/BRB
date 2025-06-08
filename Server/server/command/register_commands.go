package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

func RegisterCommands() {
	// Dev Commands
	cmd.Register(cmd.New("debug", "To debug certain code", nil, DebugCommand{}))

	// Admin Commands
	cmd.Register(cmd.New("color", "To show list of colors", nil, ColorCommand{}))
	cmd.Register(cmd.New("rank", "To set a specific player's rank", nil, RankCommand{}))

	// Player Commands
	cmd.Register(cmd.New("ping", "To check a specific player's ping", nil, PingCommand{}))
	cmd.Register(cmd.New("link", "To link your minecraft account with your discord account", nil, LinkCommand{}))
	cmd.Register(cmd.New("join", "To join a ranked BedWars game you queued in from Discord", nil, JoinCommand{}))
	cmd.Register(cmd.New("warp", "To teleport all online IDLE players from your ranked BedWars queue to your game", nil, WarpCommand{}))

	// Game Commands
	cmd.Register(cmd.New("hub", "To teleport back to hub", nil, HubCommand{}))
}
