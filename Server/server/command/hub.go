package command

import (
	"server/server"
	"server/server/language"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type HubCommand struct{}

func (h HubCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		//if cooldowns.CanExecute(pl, cooldowns.HubCMD, false, true, 5*time.Second) {
		o.Print(text.Colourf(language.Translate(pl).Global.Commands.Success.Hub, server.Config.Prefix))
		//}
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
