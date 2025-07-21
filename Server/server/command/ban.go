package command

import (
	"fmt"
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

	fmt.Println(1)
	dt, err := server.Database.FindPlayerByName(b.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
	if err != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
		return
	}

	fmt.Println(2)
	if u.Data.Rank() > dt.Rank() {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
		return
	}

	now := time.Now()

	fmt.Println(3)
	if user.ActiveBan(dt) != nil {
		pl.Message(text.Colourf(language.Translate(pl).Commands.Error.AlreadyBanned))
		return
	}

	fmt.Println(4)
	punishment := &database.PunishmentData{
		ID:            utils.RandString(8),
		PunishedBy:    pl.Name(),
		PunishedSince: now,
		Reason:        string(b.Reason),
		RemovedBy:     "",
	}

	fmt.Println(5)
	var dur time.Duration
	if b.Duration == "permanent" {
		punishment.Permanent = true
	} else {
		dur = b.Duration.Parse()
		punishment.EndsAt = now.Add(dur)
	}

	fmt.Println(6)
	dt.Punishments.Bans = append(dt.Punishments.Bans, punishment)

	durationStr := lo.If(punishment.Permanent, "").Else(" for " + utils.FriendlyDuration(dur))

	fmt.Println(7)
	if h, ok := server.MCServer.Player(dt.UUID); ok {
		go h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			dt.Punishments.Ban(e.(*player.Player), punishment)
		})
	}

	fmt.Println(8)
	user.UpdateUserData(dt)
	fmt.Println(9)
	utils.Panic(server.Database.SavePlayer(dt))

	fmt.Println(10)
	pl.Message(text.Colourf(
		language.Translate(pl).Commands.Success.Ban,
		server.Config.Prefix,
		dt.Username,
		durationStr,
		punishment.Reason,
	))
}
