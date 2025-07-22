package command

import (
	"fmt"
	"math/rand"
	"server/server"
	"server/server/api"
	"server/server/language"
	"server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type LinkCommand struct{}

func (LinkCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u.IsCooldownActive(user.CommandLink, 30*time.Second, false, true, true) {
			return
		}

		api.RegistrationCodes[pl.UUID()] = fmt.Sprintf("%04d", rand.Intn(10000))
		time.AfterFunc(10*time.Minute, func() {
			if api.RegistrationCodes[pl.UUID()] != "" {
				delete(api.RegistrationCodes, pl.UUID())
				pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LinkExpired))
			}
		})
		o.Print(text.Colourf(language.Translate(pl).Commands.Success.Link, server.Config.Prefix, api.RegistrationCodes[pl.UUID()]))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
