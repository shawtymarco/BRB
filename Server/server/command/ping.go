package command

import (
	"server/server"
	"server/server/language"
	"server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PingCommand struct {
	Target cmd.Optional[[]cmd.Target] `cmd:"target"`
}

func (p PingCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		u := user.GetUser(pl)
		if u.IsCooldownActive(user.CommandPing, 5*time.Second, false, true, true) {
			return
		}

		var target cmd.Target
		if targets, ok := p.Target.Load(); ok {
			if len(targets) != 1 {
				o.Error(text.Colourf(language.Translate(pl).Commands.Error.OnlyOneTarget))
				return
			}
			target = targets[0]
		} else {
			target = src
		}
		pl := target.(*player.Player)
		if src != target {
			targetStr := text.Colourf("<yellow>%v's</yellow>", pl.Name())
			o.Print(text.Colourf(language.Translate(pl).Commands.Success.Ping, server.Config.Prefix, targetStr, (pl.Latency() * 2).Round(time.Millisecond).Milliseconds()))
		} else {
			o.Print(text.Colourf(language.Translate(pl).Commands.Success.PingSelf, server.Config.Prefix, (pl.Latency() * 2).Round(time.Millisecond).Milliseconds()))
		}
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}
