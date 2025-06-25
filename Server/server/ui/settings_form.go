package ui

import (
	"server/server/games/bedwars"
	"time"

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
	h := submitter.(*player.Player).H()
	time.AfterFunc(700*time.Millisecond, func() {
		h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			pl := e.(*player.Player)
			switch button.Text {
			case "QuickBuy Config":
				bedwars.SendItemShopUI(&bedwars.ItemShop{Player: pl}, true)
			case "HotBar Config":
				sendHotBarConfigUI(pl)
			}
		})
	})
}

func (g SettingsForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewSettingsForm(), text.Colourf("<emerald>Settings</emerald>"))
	f = f.WithButtons(
		form.NewButton("QuickBuy Config", ""),
		form.NewButton("HotBar Config", ""),
	)
	pl.SendForm(f)
}
