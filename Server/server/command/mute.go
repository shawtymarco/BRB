package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"time"

	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type MuteCommand struct {
	Player   string      `cmd:"player"`
	Duration Duration    `cmd:"duration"`
	Reason   cmd.Varargs `cmd:"reason"`
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

	dt, err := server.Database.FindPlayerByName(m.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	if u.Data.Rank() > dt.Rank() {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
		return
	}

	now := time.Now()

	if user.ActiveMute(dt) != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.AlreadyMuted))
		return
	}

	punishment := &database.PunishmentData{
		PunishedBy:    pl.Name(),
		PunishedSince: now,
		Reason:        string(m.Reason),
		RemovedBy:     "",
	}

	var dur time.Duration
	if m.Duration == "permanent" {
		punishment.Permanent = true
	} else {
		dur = m.Duration.Parse()
		punishment.EndsAt = now.Add(dur)
	}

	dt.Punishments.Mutes = append(dt.Punishments.Mutes, punishment)

	durationStr := lo.If(punishment.Permanent, "").Else("for " + utils.FriendlyDuration(dur))

	if h, ok := server.MCServer.Player(dt.UUID); ok {
		go h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			e.(*player.Player).Message(text.Colourf(
				language.Translate(e.(*player.Player)).Commands.Success.Muted,
				lo.If(punishment.Permanent, "permanently").Else("temporarily"),
				durationStr,
				punishment.Reason,
			))
		})
	}

	user.UpdateUserData(dt)
	utils.Panic(server.Database.SavePlayer(dt))

	pl.Message(text.Colourf(
		language.Translate(pl).Commands.Success.Mute,
		server.Config.Prefix,
		dt.Username,
		durationStr,
		punishment.Reason,
	))
}
