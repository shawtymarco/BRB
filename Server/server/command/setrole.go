package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SetRoleCommand struct {
	Targets []cmd.Target `cmd:"target"`
	Rank    Rank         `cmd:"role"`
}

func (SetRoleCommand) Allow(src cmd.Source) bool {
	return SetRole.Test(src)
}

func (SetRoleCommand) PermissionMessage(src cmd.Source) string {
	return SetRole.PermissionMessage(src)
}

func (r SetRoleCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		if len(r.Targets) != 1 {
			o.Error(text.Colourf(language.Translate(pl).Commands.Error.OnlyOneTarget))
			return
		}
		t := r.Targets[0].(*player.Player)
		oldName := t.NameTag()
		r := database.RankFromPrefix(string(r.Rank))
		ds := utils.Panics(server.Database.FindPlayer(pl.UUID()))
		dt := utils.Panics(server.Database.FindPlayer(t.UUID()))
		if ds.Statistics.RankId >= r.Shortened() {
			pl.Message(text.Colourf(language.Translate(pl).Commands.Error.RankHierarchy))
			return
		}
		dt.Statistics.RankId = r.Shortened()
		utils.Panic(server.Database.SavePlayer(dt))
		t.SetNameTag(database.LobbyNameDisplay.Name(dt))
		pl.Message(text.Colourf(language.Translate(pl).Commands.Success.GiveRank, server.Config.Prefix, oldName, r.Prefix()))
	} else {
		o.Error(text.Colourf("<red>You cannot use this command in console. Please execute it in-game.</red>"))
	}
}

type Rank string

func (Rank) Type() string {
	return "rank"
}

func (Rank) Options(_ cmd.Source) []string {
	var coloredRankPrefixes []string
	for _, rp := range database.RankPrefixes {
		rank := database.RankFromPrefix(rp)
		coloredRankPrefixes = append(coloredRankPrefixes, text.Colourf(rank.Prefix()))
	}
	return coloredRankPrefixes
}
