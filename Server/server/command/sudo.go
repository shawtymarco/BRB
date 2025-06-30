package command

import (
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type SudoCommand struct {
	Player  ArgumentPlayer `cmd:"player"`
	Command cmd.Varargs    `cmd:"command"`
}

func (SudoCommand) Allow(src cmd.Source) bool {
	return Sudo.Test(src)
}

func (SudoCommand) PermissionMessage(src cmd.Source) string {
	return Sudo.PermissionMessage(src)
}

func (r SudoCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	utils.Panic(r.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		target.ExecuteCommand(string(r.Command))
	}))
}
