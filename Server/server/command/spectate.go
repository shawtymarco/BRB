package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SpectateCommand struct {
	Player ArgumentPlayer `cmd:"player"`
}

func (SpectateCommand) Allow(src cmd.Source) bool {
	return Spectate.Test(src)
}

func (SpectateCommand) PermissionMessage(src cmd.Source) string {
	return Spectate.PermissionMessage(src)
}

func (r SpectateCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)

	if pl.Name() == string(r.Player) {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.CannotSpectateOneSelf))
		return
	}
	if u.Game != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
		return
	}

	utils.Panic(r.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		ut := user.GetUser(target)

		if ut.Game == nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.MustBeInGame))
			return
		}

		u.Game = ut.Game

		pl.SetGameMode(world.GameModeSpectator)
		pl.Inventory().Clear()
		pl.Armour().Clear()
		tx.RemoveEntity(pl)
		tgtTx.AddEntity(pl.H())
		pl.Teleport(target.Position())

		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.Spectate, server.Config.Prefix, database.LobbyNameDisplay.Name(ut.Data)))
	}))
}
