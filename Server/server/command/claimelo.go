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

type ClaimELOCommand struct {
}

func (ClaimELOCommand) Allow(src cmd.Source) bool {
	return ClaimELO.Test(src)
}

func (ClaimELOCommand) PermissionMessage(src cmd.Source) string {
	return ClaimELO.PermissionMessage(src)
}

func (ClaimELOCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u.Data.Cosmetics.ELOClaimed {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.ELOAlreadyClaimed))
			return
		}

		u.Data.Cosmetics.ELOClaimed = true
		u.Data.Statistics.ELO += 250

		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.ELOClaim, server.Config.Prefix))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
