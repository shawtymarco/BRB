package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type DebugCommand struct{}

func (d DebugCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if _, ok := src.(*player.Player); ok {
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
