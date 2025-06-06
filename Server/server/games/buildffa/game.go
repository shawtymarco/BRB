package buildffa

import (
	"server/server"
	"server/server/game"
	"server/server/user"
	"server/server/utils"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/player"
)

var Game *BuildFFA

type BuildFFA struct {
	*game.Game
}

func NewBuildFFA() {
	Game = &BuildFFA{Game: game.NewGame(uuid.New(), utils.Panics(server.WorldManager.World(game.TypeBuildFFA.Title())), "grey")}
	go func() {
		for range time.NewTicker(250 * time.Millisecond).C {
			var users []*user.User
			Game.ForEachPlayer(func(pl *player.Player) {
				users = append(users, user.LookupPlayer(pl))
			})

			slices.SortFunc(users, func(a, b *user.User) int {
				return b.GameInfo.BuildFFA.Kills - a.GameInfo.BuildFFA.Kills
			})

			Game.ForEachPlayer(func(pl *player.Player) {
				u := user.LookupPlayer(pl)
				u.Scoreboard.Set(0, text.Colourf("         <yellow>▷ <white>Season 1</white> ◁</yellow>"))
				u.Scoreboard.Set(1, "§0 ")
				u.Scoreboard.Set(2, text.Colourf("<grey>Game:</grey> <emerald>Build UHC</emerald>"))
				u.Scoreboard.Set(3, "§1 ")
				u.Scoreboard.Set(4, text.Colourf("<grey>Leaderboard</grey>"))
				if len(users) > 0 {
					u.Scoreboard.Set(5, text.Colourf("<grey>1. %v</grey> <black>-</black> <emerald>%v</emerald>", users[0].Data.Username, users[0].GameInfo.BuildFFA.Kills))
				}
				if len(users) > 1 {
					u.Scoreboard.Set(6, text.Colourf("<grey>2. %v</grey> <black>-</black> <emerald>%v</emerald>", users[1].Data.Username, users[1].GameInfo.BuildFFA.Kills))
				}
				if len(users) > 2 {
					u.Scoreboard.Set(7, text.Colourf("<grey>3. %v</grey> <black>-</black> <emerald>%v</emerald>", users[2].Data.Username, users[2].GameInfo.BuildFFA.Kills))
				}
				u.Scoreboard.Set(8, "§2 ")
				u.Scoreboard.Set(9, text.Colourf("<grey>Kills:</grey> <emerald>%v</emerald>", u.GameInfo.BuildFFA.Kills))
				u.Scoreboard.Set(10, text.Colourf("<grey>Position:</grey> <emerald>#%v</emerald>", slices.Index(users, u)+1))
				u.Scoreboard.Set(11, "§3 ")
				u.Scoreboard.Set(12, text.Colourf("<yellow>ELIAGIC.CLUB</yellow>"))
				pl.SendScoreboard(u.Scoreboard)
			})
		}
	}()
}

func (b *BuildFFA) Type() game.TypeGame {
	return game.TypeBuildFFA
}

func (b *BuildFFA) Maps() []string {
	return []string{"./maps/BuildFFA"}
}

func (b *BuildFFA) MapConfig() game.MapData {
	return game.Maps[b.Maps()[0]]
}

func (b *BuildFFA) Handler() player.Handler {
	return Handler{}
}

func (b *BuildFFA) Reward(player *player.Player) {
}
