package command

import (
	"server/server"
	"server/server/language"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type GamemodeCommand struct {
	Targets  []cmd.Target `cmd:"target"`
	GameMode GameModeType `cmd:"gamemode"`
}

func (GamemodeCommand) Allow(src cmd.Source) bool {
	return Gamemode.Test(src)
}

func (GamemodeCommand) PermissionMessage(src cmd.Source) string {
	return Gamemode.PermissionMessage(src)
}

func (g GamemodeCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		if len(g.Targets) != 1 {
			o.Error(text.Colourf(language.Translate(pl).Global.Commands.Error.OnlyOneTarget))
			return
		}
		t := g.Targets[0].(*player.Player)
		switch g.GameMode {
		case "survival":
			t.SetGameMode(world.GameModeSurvival)
		case "creative":
			t.SetGameMode(world.GameModeCreative)
		case "adventure":
			t.SetGameMode(world.GameModeAdventure)
		case "spectator":
			t.SetGameMode(world.GameModeSpectator)
		}
		pl.Message(text.Colourf(language.Translate(pl).Global.Commands.Success.GameMode, server.Config.Prefix, strings.ToTitle(string(g.GameMode))))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type GameModeType string

func (GameModeType) Type() string {
	return "game mode"
}

func (GameModeType) Options(_ cmd.Source) []string {
	return []string{"survival", "creative", "adventure", "spectator"}
}
