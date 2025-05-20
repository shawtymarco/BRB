package command

import (
	"server/server"
	"server/server/language"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type GameStartCommand struct {
	Start cmd.SubCommand `cmd:"start"`
	Count int            `cmd:"timer"`
}

func (GameStartCommand) Allow(src cmd.Source) bool {
	return GameForceStart.Test(src)
}

func (GameStartCommand) PermissionMessage(src cmd.Source) string {
	return GameForceStart.PermissionMessage(src)
}

func (gsc GameStartCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		pl.Message(text.Colourf(language.Translate(pl).Global.Game.ForceStartGame, server.Config.Prefix))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
