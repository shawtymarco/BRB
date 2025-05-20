package command

import (
	"server/server"
	"server/server/database"
	"server/server/language"
	"server/server/user"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type RankCommand struct {
	Targets []cmd.Target `cmd:"target"`
	Rank    Rank         `cmd:"rank"`
}

func (RankCommand) Allow(src cmd.Source) bool {
	return GiveRank.Test(src)
}

func (RankCommand) PermissionMessage(src cmd.Source) string {
	return GiveRank.PermissionMessage(src)
}

func (r RankCommand) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if pl, ok := src.(*player.Player); ok {
		if len(r.Targets) != 1 {
			o.Error(text.Colourf(language.Translate(pl).Global.Commands.Error.OnlyOneTarget))
			return
		}
		t := r.Targets[0].(*player.Player)
		oldName := t.NameTag()
		r := database.RankFromPrefix(string(r.Rank))
		if user.DataFromPlayer(pl).GroupSettings.RankId >= r.Shortened() {
			pl.Message(text.Colourf(language.Translate(pl).Global.Commands.Error.RankHierarchy))
			return
		}
		user.DataFromPlayer(t).GroupSettings.RankId = r.Shortened()
		pl.SetNameTag(database.BasicNameDisplay.Name(user.DataFromPlayer(pl)))
		pl.Message(text.Colourf(language.Translate(pl).Global.Commands.Success.GiveRank, server.Config.Prefix, oldName, r.Prefix()))
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
	for _, r := range database.RankPrefixes {
		coloredRankPrefixes = append(coloredRankPrefixes, text.Colourf(r))
	}
	return coloredRankPrefixes
}
