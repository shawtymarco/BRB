package command

import (
	"server/server"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type GameModeCommand struct {
	Mode gameMode `cmd:"mode"`
}

func (GameModeCommand) Allow(src cmd.Source) bool {
	return GameMode.Test(src)
}

func (GameModeCommand) PermissionMessage(src cmd.Source) string {
	return GameMode.PermissionMessage(src)
}

func (c GameModeCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u.Game != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
			return
		}

		var mode world.GameMode
		var modeName string

		switch string(c.Mode) {
		case "survival", "s", "0":
			mode = world.GameModeSurvival
			modeName = "Survival"
		case "creative", "c", "1":
			mode = world.GameModeCreative
			modeName = "Creative"
		case "adventure", "a", "2":
			mode = world.GameModeAdventure
			modeName = "Adventure"
		case "spectator", "sp", "3":
			mode = world.GameModeSpectator
			modeName = "Spectator"
		default:
			pl.Message(text.Colourf("<red>Invalid game mode. Use: survival, creative, adventure, or spectator</red>"))
			return
		}

		pl.SetGameMode(mode)
		pl.Message(text.Colourf("<green>%sGame mode set to %s</green>", server.Config.Prefix, modeName))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type gameMode string

func (gameMode) Type() string {
	return "gameMode"
}

func (gameMode) Options(src cmd.Source) []string {
	return []string{"survival", "creative", "adventure", "spectator", "s", "c", "a", "sp", "0", "1", "2", "3"}
}
