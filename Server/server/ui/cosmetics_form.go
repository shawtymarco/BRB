package ui

import (
	"server/server/database"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type CosmeticsForm struct {
}

func NewCosmeticsForm() CosmeticsForm {
	return CosmeticsForm{}
}

func (g CosmeticsForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	switch button.Text {
	case "Wood Skins":
		NewWoodSkinsForm().SendTo(pl)
		break
	case "Capes":
		NewCapesForm().SendTo(pl)
		break
	}
}

func (g CosmeticsForm) Close(submitter form.Submitter) {
	g.SendTo(submitter.(*player.Player))
}

func (g CosmeticsForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewCosmeticsForm(), text.Colourf("<emerald>Cosmetics</emerald>"))
	f = f.WithButtons(
		form.NewButton("Wood Skins", ""),
		form.NewButton("Capes", ""),
	)
	pl.SendForm(f)
}

type WoodSkinsForm struct {
}

func NewWoodSkinsForm() WoodSkinsForm {
	return WoodSkinsForm{}
}

func (w WoodSkinsForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	u := user.LookupPlayer(pl)

	wt := Load[block.WoodType](pl, button.Text)
	u.Data.Cosmetics.SelectedWoodType = wt

	pl.Message(text.Colourf(language.Translate(pl).Misc.SelectedWoodSkin, wt.Name()))
}

func (w WoodSkinsForm) Close(submitter form.Submitter) {
	w.SendTo(submitter.(*player.Player))
}

func (w WoodSkinsForm) SendTo(pl *player.Player) {
	u := user.LookupPlayer(pl)
	f := form.NewMenu(NewWoodSkinsForm(), text.Colourf("<emerald>Wood Skins</emerald>"))
	for _, wt := range block.WoodTypes() {
		f = f.WithButtons(
			AddButtonWithValue(
				pl,
				text.Colourf("%v%v", wt.Name(), lo.If(u.Data.Cosmetics.SelectedWoodType == wt, text.Colourf("\n<green>Selected</green>")).Else("")),
				"",
				wt,
			),
		)
	}
	pl.SendForm(f)
}

type CapesForm struct {
}

func NewCapesForm() CapesForm {
	return CapesForm{}
}

func (c CapesForm) Submit(submitter form.Submitter, button form.Button, tx *world.Tx) {
	pl := submitter.(*player.Player)
	u := user.LookupPlayer(pl)
	cape := Load[database.Cape](pl, button.Text)
	u.Data.Cosmetics.SelectedCape = cape.Identifier()
	u.RefreshCape()

	pl.Message(text.Colourf(language.Translate(pl).Misc.SelectedCape, cape.Name()))
}

func (c CapesForm) Close(submitter form.Submitter) {
	c.SendTo(submitter.(*player.Player))
}

func (c CapesForm) SendTo(pl *player.Player) {
	u := user.LookupPlayer(pl)
	f := form.NewMenu(NewCapesForm(), text.Colourf("<emerald>Capes</emerald>"))
	for _, cape := range database.AllCapes() {
		f = f.WithButtons(
			AddButtonWithValue(
				pl,
				text.Colourf("%v%v", cape.Name(), lo.If(u.Data.Cosmetics.SelectedCape == cape.Identifier(), text.Colourf("\n<green>Selected</green>")).Else("")),
				"",
				cape,
			),
		)
	}
	pl.SendForm(f)
}
