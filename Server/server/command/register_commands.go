package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

func RegisterCommands() {

	// Global Commands
	cmd.Register(cmd.New("hub", "To teleport back to hub", nil, HubCommand{}))
	cmd.Register(cmd.New("ping", "To check a specific player's ping", nil, PingCommand{}))
	cmd.Register(cmd.New("whisper", "To send a private message to the specified player", []string{"w", "tell", "msg"}, WhisperCommand{}))
	cmd.Register(cmd.New("reply", "To reply privately to the last player that messaged you", []string{"r"}, ReplyCommand{}))
	cmd.Register(cmd.New("link", "To link your minecraft account with your discord account", nil, LinkCommand{}))
	cmd.Register(cmd.New("join", "To join a RBW game you queued in from Discord", nil, JoinCommand{}))
	cmd.Register(cmd.New("rejoin", "To rejoin the last game you joined", nil, JoinCommand{}))
	cmd.Register(cmd.New("warp", "To teleport all online IDLE players from your RBW queue to your game", nil, WarpCommand{}))
	cmd.Register(cmd.New("duel", "To duel another player", []string{"duels"}, DuelRequestCommand{}, DuelAcceptCommand{}))

	// Perk Commands
	cmd.Register(cmd.New("fly", "To fly in the lobby", nil, FlyCommand{}))
	cmd.Register(cmd.New("spectate", "To spectate a specific player", nil, SpectateCommand{}))
	cmd.Register(cmd.New("nick", "To change your nickname", nil, NickCommand{}))
	cmd.Register(cmd.New("claimelo", "To claim your seasonal 250 ELO", nil, ClaimELOCommand{}))

	// Staff Commands
	cmd.Register(cmd.New("ban", "To ban a player", nil, BanCommand{}))
	cmd.Register(cmd.New("unban", "To unban a player", nil, UnbanCommand{}))
	cmd.Register(cmd.New("mute", "To mute a player", nil, MuteCommand{}))
	cmd.Register(cmd.New("unmute", "To unmute a player", nil, UnmuteCommand{}))
	cmd.Register(cmd.New("alias", "To see all alt accounts of a player", nil, AliasCommand{}))

	// Admin Commands
	cmd.Register(cmd.New("setrole", "To set a specific player's rank", nil, SetRoleCommand{}))
	cmd.Register(cmd.New("addcape", "To add a cape to a specific player", nil, AddCapeCommand{}))
	cmd.Register(cmd.New("removecape", "To remove a cape from a specific player", nil, RemoveCapeCommand{}))
	cmd.Register(cmd.New("resetstats", "To reset a specific player's statistics", nil, ResetStatsCommand{}))
	cmd.Register(cmd.New("gamemode", "To change your game mode", []string{"gm"}, GameModeCommand{}))
	//cmd.Register(cmd.New("sudo", "To execute commands from another player's perspective", nil, SudoCommand{}))
}
