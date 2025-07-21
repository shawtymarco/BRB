package ui

import (
	"fmt"
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
	fmt.Println(pl.Name())
	fmt.Println(1)
	gt := Load[game.TypeGame](pl, button.Text)
	pl.Handler().HandleQuit(pl)

	fmt.Println(2)
	switch gt {
	case game.TypeBuildFFA:
		fmt.Println(3)
		buildffa.Join(pl, tx)
		fmt.Println(4)
	case game.TypeBedFight:
		fmt.Println(5)
		bedwars.Join(pl, tx, 1, 2, game.TypeBedFight, false, nil)
		fmt.Println(6)
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
