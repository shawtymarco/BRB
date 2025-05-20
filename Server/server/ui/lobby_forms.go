package ui

import (
	"server/server"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

type WelcomeForm struct {
	Play form.Button
	//Store     form.Button
}

func NewWelcomeForm() WelcomeForm {
	return WelcomeForm{
		Play: form.Button{
			Text: "Play",
		},
		//Store: form.Button{
		//	Text: "Store",
		//},
	}
}

func (w WelcomeForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	if button == w.Play {
		GamesForm{}.SendTo(pl)
	}
}

func (w WelcomeForm) Close(submitter form.Submitter) {
	w.SendTo(submitter.(*player.Player))
}

func (w WelcomeForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewWelcomeForm(), "home.menu")
	pl.SendForm(f)
}

type GamesForm struct {
	isBuilder bool
}

func NewGamesForm() GamesForm {
	return GamesForm{}
}

func (g GamesForm) Submit(submitter form.Submitter, button form.Button, _ *world.Tx) {
	pl := submitter.(*player.Player)
	sid := Load[server.SrvId](pl, button.Text)

	u := user.LookupPlayer(pl)
	u.Transfer(sid)
}

func (g GamesForm) Close(submitter form.Submitter) {
	g.SendTo(submitter.(*player.Player))
}

func (g GamesForm) SendTo(pl *player.Player) {
	f := form.NewMenu(NewGamesForm(), GameSelector)
	f = f.WithButtons(
		AddButtonWithValue(pl, server.TowerWars.DisplayName(), "textures/assets/tower_defense_v2.png", server.TowerWars),
		AddButtonWithValue(pl, server.Dungeons.DisplayName(), "", server.Dungeons),
	)
	pl.SendForm(f)
}
