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

type UnbanCommand struct {
	Player string `cmd:"player"`
}

func (UnbanCommand) Allow(src cmd.Source) bool {
	return Ban.Test(src)
}

func (UnbanCommand) PermissionMessage(src cmd.Source) string {
	return Ban.PermissionMessage(src)
}

func (b UnbanCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}

	dt, err := server.Database.FindPlayerByName(b.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	activeBan := user.ActiveBan(dt)
	if activeBan == nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NotBanned))
		return
	}

	activeBan.RemovedBy = pl.Name()

	user.UpdateUserData(dt)
	utils.Panic(server.Database.SavePlayer(dt))

	pl.Message(text.Colourf(language.Translate(pl).Commands.Success.Unban, server.Config.Prefix))
}
