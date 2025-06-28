package command

import (
	"server/server"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type FlyCommand struct {
}

func (FlyCommand) Allow(src cmd.Source) bool {
	return Fly.Test(src)
}

func (FlyCommand) PermissionMessage(src cmd.Source) string {
	return Fly.PermissionMessage(src)
}

func (FlyCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u.Game != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.LobbyOnly))
			return
		}

		if pl.GameMode() == (GameModeFly{}) {
			pl.StopFlying()
			pl.SetGameMode(world.GameModeSurvival)
			pl.Message(text.Colourf(language.Translate(pl).Commands.Success.FlyOff, server.Config.Prefix))
		} else {
			pl.SetGameMode(GameModeFly{})
			pl.StartFlying()
			pl.Message(text.Colourf(language.Translate(pl).Commands.Success.FlyOn, server.Config.Prefix))
		}
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type GameModeFly struct {
}

func (g GameModeFly) AllowsEditing() bool {
	return false
}

func (g GameModeFly) AllowsTakingDamage() bool {
	return false
}

func (g GameModeFly) CreativeInventory() bool {
	return false
}

func (g GameModeFly) HasCollision() bool {
	return true
}

func (g GameModeFly) AllowsFlying() bool {
	return true
}

func (g GameModeFly) AllowsInteraction() bool {
	return true
}

func (g GameModeFly) Visible() bool {
	return true
}
