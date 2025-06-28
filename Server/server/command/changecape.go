package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"
	"slices"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type AddCapeCommand struct {
	Player ArgumentPlayer `cmd:"player"`
	Cape   cape           `cmd:"cape"`
}

func (AddCapeCommand) Allow(src cmd.Source) bool {
	return ChangeCape.Test(src)
}

func (AddCapeCommand) PermissionMessage(src cmd.Source) string {
	return ChangeCape.PermissionMessage(src)
}

func (c AddCapeCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if slices.Contains(u.Data.Cosmetics.OwnedCapes, string(c.Cape)) {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.CapeAlreadyOwned))
			return
		}

		u.Data.Cosmetics.OwnedCapes = append(u.Data.Cosmetics.OwnedCapes, string(c.Cape))
		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.AddCape, server.Config.Prefix))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type RemoveCapeCommand struct {
	Player ArgumentPlayer `cmd:"player"`
	Cape   cape           `cmd:"cape"`
}

func (RemoveCapeCommand) Allow(src cmd.Source) bool {
	return ChangeCape.Test(src)
}

func (RemoveCapeCommand) PermissionMessage(src cmd.Source) string {
	return ChangeCape.PermissionMessage(src)
}

func (c RemoveCapeCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if !slices.Contains(u.Data.Cosmetics.OwnedCapes, string(c.Cape)) {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.CapeNotOwned))
			return
		}

		for i, cp := range u.Data.Cosmetics.OwnedCapes {
			if cp == string(c.Cape) {
				u.Data.Cosmetics.OwnedCapes = append(u.Data.Cosmetics.OwnedCapes[:i], u.Data.Cosmetics.OwnedCapes[i+1:]...)
				break
			}
		}

		if u.Data.Cosmetics.SelectedCape == string(c.Cape) {
			u.Data.Cosmetics.SelectedCape = ""
			u.RefreshCape()
		}

		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.RemoveCape, server.Config.Prefix))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type cape string

func (cape) Type() string {
	return "cape"
}

func (cape) Options(src cmd.Source) (capes []string) {
	for _, cape := range database.AllCapes() {
		capes = append(capes, cape.Identifier())
	}
	return capes
}
