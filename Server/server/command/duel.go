package command

import (
	"server/server"
	"server/server/database"
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/language"
	"server/server/user"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type DuelRequestCommand struct {
	Request cmd.SubCommand `cmd:"request"`
	Player  ArgumentPlayer `cmd:"player"`
}

func (d DuelRequestCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)
	if u.Game != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
		return
	}

	sendMessage := func(tx2 *world.Tx) {
		utils.Panic(d.Player.ExecWithPlayerSafe(tx2, func(tgtTx *world.Tx, target *player.Player) {
			u.PendingDuelRequestTo = target.H()
			target.Message(text.Colourf(language.Translate(target).Commands.Success.DuelRequestReceived, server.Config.Prefix, database.LobbyNameDisplay.Name(u.Data), u.Data.Username))
			pl.Message(text.Colourf(language.Translate(pl).Commands.Success.DuelRequest, server.Config.Prefix))
		}))
	}

	if pdrt := u.PendingDuelRequestTo; pdrt != nil {
		go pdrt.ExecWorld(func(tx *world.Tx, e world.Entity) {
			e.(*player.Player).Message(text.Colourf(language.Translate(e.(*player.Player)).Commands.Success.DuelRevoked, database.LobbyNameDisplay.Name(u.Data)))
			sendMessage(tx)
		})
	} else {
		sendMessage(tx)
	}
}

type DuelAcceptCommand struct {
	Accept cmd.SubCommand `cmd:"accept"`
	Player ArgumentPlayer `cmd:"player"`
}

func (d DuelAcceptCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)
	if u.Game != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
		return
	}

	utils.Panic(d.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		ut := user.GetUser(target)
		if ut.PendingDuelRequestTo == nil || ut.PendingDuelRequestTo.UUID() != pl.UUID() {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Success.NoDuelRequest))
			return
		}

		ut.PendingDuelRequestTo = nil

		b := bedwars.NewBedWars(game.TypeBedFight, 1, 2, false)
		bedwars.Join(pl, tx, 1, 2, game.TypeBedFight, false, b)
		bedwars.Join(target, tgtTx, 1, 2, game.TypeBedFight, false, b)
	}))
}
