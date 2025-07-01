package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type AliasCommand struct {
	Player string `cmd:"player"`
}

func (AliasCommand) Allow(src cmd.Source) bool {
	return Alias.Test(src)
}

func (AliasCommand) PermissionMessage(src cmd.Source) string {
	return Alias.PermissionMessage(src)
}

func (a AliasCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}

	dt, err := server.Database.FindPlayerByName(a.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	var usernames []string
	for _, acc := range dt.AlternativeMCAccounts {
		usernames = append(usernames, acc.Username)
	}

	pl.Message(text.Colourf(language.Translate(pl).Commands.Success.Alias, server.Config.Prefix, strings.Join(usernames, ", ")))
}
