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

type MuteCommand struct {
	Player   ArgumentPlayer `cmd:"player"`
	Duration Duration       `cmd:"duration"`
	Reason   string         `cmd:"reason"`
}

func (MuteCommand) Allow(src cmd.Source) bool {
	return Mute.Test(src)
}

func (MuteCommand) PermissionMessage(src cmd.Source) string {
	return Mute.PermissionMessage(src)
}

func (m MuteCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	pl, ok := src.(*player.Player)
	if !ok {
		o.Error(text.Colourf("<red>You must run this command in-game.</red>"))
		return
	}
	u := user.GetUser(pl)

	utils.Panic(m.Player.ExecWithPlayerSafe(tx, func(tgtTx *world.Tx, target *player.Player) {
		ut := user.GetUser(target)
		if u.Data.Rank() > ut.Data.Rank() {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
			return
		}

		now := time.Now()

		if u.Data.Punishments.ActiveMute() != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.AlreadyMuted))
			return
		}

		punishment := &database.PunishmentData{
			PunishedBy:    pl.Name(),
			PunishedSince: now,
			Reason:        m.Reason,
			RemovedBy:     "",
		}

		var dur time.Duration
		if m.Duration == "permanent" {
			punishment.Permanent = true
		} else {
			dur = m.Duration.Parse()
			punishment.EndsAt = now.Add(dur)
		}

		u.Data.Punishments.Mutes = append(u.Data.Punishments.Mutes, punishment)

		durationStr := lo.If(punishment.Permanent, "").Else("for " + utils.FriendlyDuration(dur))

		target.Message(text.Colourf(
			language.Translate(target).Commands.Success.Muted,
			lo.If(punishment.Permanent, "permanently").Else("temporarily"),
			durationStr,
			punishment.Reason,
		))

		pl.Message(text.Colourf(
			language.Translate(pl).Commands.Success.Mute,
			server.Config.Prefix,
			target.Name(),
			durationStr,
			punishment.Reason,
		))
	}))
}
