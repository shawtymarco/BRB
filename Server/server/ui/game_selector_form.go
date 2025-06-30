package ui

import (
	"server/server"
	"server/server/game"
	"server/server/games/bedwars"
	"server/server/games/buildffa"
	"time"

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

func (g GameSelectorForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	gt := Load[game.TypeGame](pl, button.Text)
	pl.Handler().HandleQuit(pl)

	switch gt {
	case game.TypeBuildFFA:
		buildffa.Join(pl, pl.Tx())

		if server.BotMark != nil {
			time.AfterFunc(5*time.Second, func() {
				server.BotMark.H().ExecWorld(func(tx2 *world.Tx, e world.Entity) {
					buildffa.Join(e.(*player.Player), tx2)
				})
			})
		}
	case game.TypeBedFight:
		bedwars.Join(pl, pl.Tx(), 1, 2, game.TypeBedFight, false, nil)
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
