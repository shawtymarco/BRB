package command

import (
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

func (JoinCommand) Allow(src cmd.Source) bool {
	return Join.Test(src)
}

func (JoinCommand) PermissionMessage(src cmd.Source) string {
	return GiveRank.PermissionMessage(src)
}

func (r JoinCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.LookupPlayer(pl)
		for _, g := range bedwars.Games {
			if slices.Contains(g.UsersToJoin, u.Data.UserId) {
				bedwars.Join(pl, pl.Tx(), g.TeamSize, g.TeamCount, g.Type(), false, g)
				return
			}
		}

		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NoGameToJoin))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
