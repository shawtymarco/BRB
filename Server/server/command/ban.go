package command

import (
	"github.com/samber/lo"
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type BanCommand struct {
	Player   ArgumentPlayer `cmd:"player"`
	Duration Duration       `cmd:"duration"`
	Reason   string         `cmd:"reason"`
}

func (BanCommand) Allow(src cmd.Source) bool {
	return Ban.Test(src)
}

func (BanCommand) PermissionMessage(src cmd.Source) string {
	return Ban.PermissionMessage(src)
}

func (b BanCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)

	utils.Panic(b.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		ut := user.GetUser(target)
		if u.Data.Rank() > ut.Data.Rank() {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
			return
		}

		now := time.Now()

		if u.Data.Punishments.ActiveBan() != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.AlreadyBanned))
			return
		}

		punishment := &database.PunishmentData{
			PunishedBy:    pl.Name(),
			PunishedSince: now,
			Reason:        b.Reason,
			RemovedBy:     "",
		}

		var dur time.Duration
		if b.Duration == "permanent" {
			punishment.Permanent = true
		} else {
			dur = b.Duration.Parse()
			punishment.EndsAt = now.Add(dur)
		}

		u.Data.Punishments.Bans = append(u.Data.Punishments.Bans, punishment)

		durationStr := lo.If(punishment.Permanent, "").Else("for " + utils.FriendlyDuration(dur))

		target.Disconnect(
			language.Translate(target).Commands.Success.BanDisconnect,
			lo.If(punishment.Permanent, "permanently").Else("temporarily"),
			durationStr,
			punishment.Reason,
		)

		pl.Message(text.Colourf(
			language.Translate(pl).Commands.Success.Ban,
			server.Config.Prefix,
			target.Name(),
			durationStr,
			punishment.Reason,
		))
	}))
}
