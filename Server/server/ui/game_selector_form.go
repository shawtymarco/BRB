package ui

import (
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/games/buildffa"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

type GameSelectorForm struct {
}

func NewGamesForm() GameSelectorForm {
	return GameSelectorForm{}
}

func (g GameSelectorForm) Submit(submitter form.Submitter, button form.Button, tx *world.Tx) {
	pl := submitter.(*player.Player)
	gt := Load[game.TypeGame](pl, button.Text)
	pl.Handler().HandleQuit(pl)

	switch gt {
	case game.TypeBuildFFA:
		buildffa.Join(pl, tx)
	case game.TypeBedFight:
		bedwars.Join(pl, tx, 1, 2, game.TypeBedFight, false, nil)
		//go server.BotMark.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
		//	bedwars.Join(e.(*player.Player), tx, 1, 2, game.TypeBedFight, false, nil)
		//})
	}
}

func (g GameSelectorForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewGamesForm(), text.Colourf("<emerald>Game Selector</emerald>"))
	f = f.WithButtons(
		AddButtonWithValue(pl, "Build FFA", "", game.TypeBuildFFA),
		AddButtonWithValue(pl, "Bed Fight", "", game.TypeBedFight),
	)
	pl.SendForm(f)
}
