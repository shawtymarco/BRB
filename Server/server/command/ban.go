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

type BanCommand struct {
	Player   string      `cmd:"player"`
	Duration Duration    `cmd:"duration"`
	Reason   cmd.Varargs `cmd:"reason"`
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

	dt, err := server.Database.FindPlayerByName(b.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	if u.Data.Rank() > dt.Rank() {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
		return
	}

	now := time.Now()

	if user.ActiveBan(dt) != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.AlreadyBanned))
		return
	}

	punishment := &database.PunishmentData{
		ID:            utils.RandString(8),
		PunishedBy:    pl.Name(),
		PunishedSince: now,
		Reason:        string(b.Reason),
		RemovedBy:     "",
	}

	var dur time.Duration
	if b.Duration == "permanent" {
		punishment.Permanent = true
	} else {
		dur = b.Duration.Parse()
		punishment.EndsAt = now.Add(dur)
	}

	dt.Punishments.Bans = append(dt.Punishments.Bans, punishment)

	durationStr := lo.If(punishment.Permanent, "").Else(" for " + utils.FriendlyDuration(dur))

	if h, ok := server.MCServer.Player(dt.UUID); ok {
		go h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			dt.Punishments.Ban(e.(*player.Player), punishment)
		})
	}

	user.UpdateUserData(dt)
	utils.Panic(server.Database.SavePlayer(dt))

	pl.Message(text.Colourf(
		language.Translate(pl).Commands.Success.Ban,
		server.Config.Prefix,
		dt.Username,
		durationStr,
		punishment.Reason,
	))
}
