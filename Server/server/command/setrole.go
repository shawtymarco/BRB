package command

import (
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

type SetRoleCommand struct {
	Player   string                 `cmd:"player"`
	Rank     Rank                   `cmd:"role"`
	Duration cmd.Optional[Duration] `cmd:"duration"`
}

func (SetRoleCommand) Allow(src cmd.Source) bool {
	return SetRole.Test(src)
}

func (SetRoleCommand) PermissionMessage(src cmd.Source) string {
	return SetRole.PermissionMessage(src)
}

func (r SetRoleCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		dt, err := server.Database.FindPlayerByName(r.Player, &database.PlayerNameSearchOpts{CaseInsensitive: true, PartialMatch: true})
		if err != nil {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.PlayerNotExist))
			return
		}

		rank := database.RankFromName(string(r.Rank))

		dt.Statistics.RankId = rank.Shortened()
		if dur, ok := r.Duration.Load(); ok && dur != "permanent" {
			dt.Statistics.RankEndsIn = time.Now().Add(dur.Parse())
		} else {
			dt.Statistics.RankEndsIn = time.Time{}
		}

		user.UpdateUserData(dt)
		utils.Panic(server.Database.SavePlayer(dt))

		if h, ok := server.MCServer.Player(dt.UUID); ok {
			go h.ExecWorld(func(tx *world.Tx, e world.Entity) {
				ut := user.GetUser(e.(*player.Player))
				e.(*player.Player).SetNameTag(database.LobbyNameDisplay.Name(ut.Data))
			})
		}

		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.GiveRank, server.Config.Prefix, dt.Username, rank.Prefix()))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type Rank string

func (Rank) Type() string {
	return "rank"
}

func (Rank) Options(_ cmd.Source) []string {
	return database.ShortenedRanks
}
