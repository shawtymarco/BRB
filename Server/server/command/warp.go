package command

import (
	"server/server"
	"server/server/games/bedwars"
	"server/server/language"
	"server/server/user"
	"slices"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type WarpCommand struct {
}

func (WarpCommand) Allow(src cmd.Source) bool {
	return Warp.Test(src)
}

func (WarpCommand) PermissionMessage(src cmd.Source) string {
	return GiveRank.PermissionMessage(src)
}

func (r WarpCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.LookupPlayer(pl)
		if u.Game == nil {
			o.Error(text.Colourf(language.Translate(pl).Game.Error.NotInAGame))
			return
		}

		var bwGame *bedwars.BedWars
		for _, g := range bedwars.Games {
			if g.ID() == u.Game.ID() {
				bwGame = g
				break
			}
		}

		if bwGame != nil {
			if len(bwGame.UsersToJoin) == 0 {
				pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NoMorePlayersToWarp))
				return
			}

			server.MCServer.World().Exec(func(tx2 *world.Tx) {
				for e := range tx2.Players() {
					p := e.(*player.Player)
					u := user.LookupPlayer(p)

					if u.Game == nil && u.Data.IsRegistered() && slices.Contains(bwGame.UsersToJoin, u.Data.UserId) {
						bedwars.Join(p, tx2, bwGame.TeamSize, bwGame.TeamCount, bwGame.Type(), false, bwGame)
						p.Message(text.Colourf(language.Translate(p).Commands.Success.YouGotWarped))
					}
				}
			})
		}
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
