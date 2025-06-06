package command

import (
	"server/server"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type HubCommand struct{}

func (HubCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.LookupPlayer(pl)
		if u.IsCooldownActive(user.CommandHub, 5*time.Second, false, true) {
			return
		}

		h := pl.H()
		tx.RemoveEntity(pl)
		server.MCServer.World().Exec(func(tx *world.Tx) {
			tx.AddEntity(pl.H())
			if e, ok := h.Entity(tx); ok {
				pl = e.(*player.Player)

				pl.Handler().HandleQuit(pl)
				lobby.Join(pl)
			}
		})

		o.Print(text.Colourf(language.Translate(pl).Commands.Success.Hub, server.Config.Prefix))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
