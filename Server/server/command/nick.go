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
	Nickname string `cmd:"nickname"`
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

		if u.Game != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
			return
		}

		if len(n.Nickname) < 2 || len(n.Nickname) > 20 {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameLength))
			return
		}

		if strings.HasPrefix(n.Nickname, " ") || strings.HasSuffix(n.Nickname, " ") {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameSpace))
			return
		}

		for _, r := range n.Nickname {
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ') {
				pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameSpecialChars))
				return
			}
		}

		if strings.Contains(n.Nickname, "  ") {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.NicknameMultipleSpaces))
			return
		}

		u.Data.Cosmetics.Nickname = n.Nickname
		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.Nick, server.Config.Prefix, n.Nickname))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
