package command

import (
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/language"
	"server/server/user"
	"slices"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type JoinCommand struct {
}

func (r JoinCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		for _, g := range bedwars.Games {
			if g.Type() == game.TypeBedWars && slices.Contains(g.UsersToJoin, u.Data.UserId) {
				pl.Handler().HandleQuit(pl)
				bedwars.Join(pl, pl.Tx(), g.TeamSize, g.TeamCount, g.Type(), false, g)
				return
			}
		}

		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NoGameToJoin))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
