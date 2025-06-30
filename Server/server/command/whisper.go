package command

import (
	"server/server/user"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type WhisperCommand struct {
	Player  ArgumentPlayer `cmd:"player"`
	Message cmd.Varargs    `cmd:"message"`
}

func (r WhisperCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)

	utils.Panic(r.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		ut := user.GetUser(target)
		u.LastMessenger = target.H()
		ut.LastMessenger = pl.H()

		target.Message(text.Colourf("<yellow><bold>FROM</bold> %v:</yellow> <grey>%v</grey>", pl.Name(), r.Message))
		pl.Message(text.Colourf("<yellow><bold>TO</bold> %v:</yellow> <grey>%v</grey>", target.Name(), r.Message))
	}))
}

type ReplyCommand struct {
	Message cmd.Varargs `cmd:"message"`
}

func (r ReplyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)

	f := func(target *player.Player) {
		target.Message(text.Colourf("<yellow><bold>FROM</bold> %v:</yellow> <grey>%v</grey>", pl.Name(), r.Message))
		pl.Message(text.Colourf("<yellow><bold>TO</bold> %v:</yellow> <grey>%v</grey>", target.Name(), r.Message))
	}

	if target, ok := u.LastMessenger.Entity(tx); ok {
		f(target.(*player.Player))
	} else {
		u.LastMessenger.ExecWorld(func(tx *world.Tx, e world.Entity) {
			f(e.(*player.Player))
		})
	}
}
