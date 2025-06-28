package command

import (
	"server/server"
	"server/server/language"
	"server/server/user"
	"strings"
	"unicode"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type NickCommand struct {
	Nickname string `json:"nickname"`
}

func (NickCommand) Allow(src cmd.Source) bool {
	return Nick.Test(src)
}

func (NickCommand) PermissionMessage(src cmd.Source) string {
	return Nick.PermissionMessage(src)
}

func (n NickCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)

		if len(n.Nickname) < 3 || len(n.Nickname) > 12 {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameLength, server.Config.Prefix))
			return
		}

		if strings.HasPrefix(n.Nickname, " ") || strings.HasSuffix(n.Nickname, " ") {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameSpace, server.Config.Prefix))
			return
		}

		for _, r := range n.Nickname {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ') {
				pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameSpecialChars, server.Config.Prefix))
				return
			}
		}

		if strings.Contains(n.Nickname, "  ") {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameMultipleSpaces, server.Config.Prefix))
			return
		}

		u.Data.Cosmetics.Nickname = n.Nickname
		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.Nick, server.Config.Prefix, n.Nickname))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
