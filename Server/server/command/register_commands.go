package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

func RegisterCommands() {

	// Global Commands
	cmd.Register(cmd.New("hub", "To teleport back to hub", nil, HubCommand{}))
	cmd.Register(cmd.New("ping", "To check a specific player's ping", nil, PingCommand{}))
	cmd.Register(cmd.New("link", "To link your minecraft account with your discord account", nil, LinkCommand{}))
	cmd.Register(cmd.New("join", "To join a RBW game you queued in from Discord", nil, JoinCommand{}))
	cmd.Register(cmd.New("rejoin", "To rejoin the last game you joined", nil, JoinCommand{}))
	cmd.Register(cmd.New("warp", "To teleport all online IDLE players from your RBW queue to your game", nil, WarpCommand{}))

	// Perk Commands
	cmd.Register(cmd.New("fly", "To fly in the lobby", nil, FlyCommand{}))
	cmd.Register(cmd.New("spectate", "To spectate a specific player", nil, SpectateCommand{}))
	cmd.Register(cmd.New("nick", "To change your nickname", nil, NickCommand{}))
	cmd.Register(cmd.New("claimelo", "To claim your seasonal 250 ELO", nil, ClaimELOCommand{}))

	// Admin Commands
	cmd.Register(cmd.New("rank", "To set a specific player's rank", nil, SetRoleCommand{}))
	cmd.Register(cmd.New("addcape", "To add a cape to a specific player", nil, AddCapeCommand{}))
	cmd.Register(cmd.New("removecape", "To remove a cape from a specific player", nil, RemoveCapeCommand{}))
	cmd.Register(cmd.New("resetstats", "To reset a specific player's statistics", nil, ResetStatsCommand{}))
}
