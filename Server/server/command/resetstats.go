package command

import (
	"server/server"
	"server/server/language"
	"server/server/user"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ResetStatsCommand struct {
	Player ArgumentPlayer `cmd:"player"`
}

func (ResetStatsCommand) Allow(src cmd.Source) bool {
	return ResetStats.Test(src)
}

func (ResetStatsCommand) PermissionMessage(src cmd.Source) string {
	return ResetStats.PermissionMessage(src)
}

func (r ResetStatsCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}

	utils.Panic(r.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		u := user.GetUser(target)
		user.ResetUser(target)
		var output cmd.Output
		HubCommand{}.Run(target, &output, tgtTx)
		u.SendScoreboard(7)
		target.Message(text.Colourf(language.Translate(target).Commands.Success.ResetStatsDisconnect))
		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.ResetStats, server.Config.Prefix))
	}))
}
