package ui

import (
	"server/server/games/bedwars"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SettingsForm struct {
}

func NewSettingsForm() SettingsForm {
	return SettingsForm{}
}

func (g SettingsForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	switch button.Text {
	case "QuickBuy Config":
		bedwars.SendItemShopUI(&bedwars.ItemShop{Player: pl, IsQuickBuy: true})
		break
	case "HotBar Config":

		break
	}
}

func (g SettingsForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewCosmeticsForm(), text.Colourf("<emerald>Settings</emerald>"))
	f = f.WithButtons(
		form.NewButton("QuickBuy Config", ""),
		form.NewButton("HotBar Config", ""),
	)
	pl.SendForm(f)
}
