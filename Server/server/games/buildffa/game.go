package buildffa

import (
	"server/server"
	"server/server/font"
	"server/server/game"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"slices"
	"time"

	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/world"

	"github.com/google/uuid"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/player"
)

var Game *BuildFFA

type BuildFFA struct {
	*game.Game
}

func NewBuildFFA() {
	gameWorld := utils.Panics(server.WorldManager.World("BFFA"))
	gameWorld.SetDifficulty(world.DifficultyNormal)
	gameWorld.StopWeatherCycle()
	gameWorld.StopRaining()
	gameWorld.StopThundering()
	gameWorld.StopTime()
	gameWorld.Handle(listener.WorldHandler{})
	Game = &BuildFFA{Game: game.NewGame(uuid.New(), gameWorld, "grey")}

	go func() {
		for range time.NewTicker(250 * time.Millisecond).C {
			Game.World().Exec(func(tx *world.Tx) {
				var users []*user.User
				Game.ForEachActivePlayer(func(pl *player.Player) {
					users = append(users, user.GetUser(pl))
				}, tx)
				slices.SortFunc(users, func(a, b *user.User) int {
					return b.GameInfo.BuildFFA.Kills - a.GameInfo.BuildFFA.Kills
				})

				Game.ForEachActivePlayer(func(pl *player.Player) {
					u := user.GetUser(pl)
					u.Scoreboard.Set(1, "§0")
					u.Scoreboard.Set(2, text.Colourf("<grey>Game:</grey> <emerald>Build FFA</emerald>"))
					u.Scoreboard.Set(3, "§1")
					u.Scoreboard.Set(4, text.Colourf("<grey>Leaderboard</grey>"))
					u.Scoreboard.Set(5, "§2")
					u.Scoreboard.Set(6, "§3")
					u.Scoreboard.Set(7, "§4")
					if len(users) > 0 {
						u.Scoreboard.Set(5, text.Colourf("<grey>1. %v</grey> <black>-</black> <emerald>%v</emerald>", lo.If(users[0].Data.Cosmetics.Nickname != "", users[0].Data.Cosmetics.Nickname).Else(users[0].Data.Username), users[0].GameInfo.BuildFFA.Kills))
					}
					if len(users) > 1 {
						u.Scoreboard.Set(6, text.Colourf("<grey>2. %v</grey> <black>-</black> <emerald>%v</emerald>", lo.If(users[1].Data.Cosmetics.Nickname != "", users[1].Data.Cosmetics.Nickname).Else(users[1].Data.Username), users[1].GameInfo.BuildFFA.Kills))
					}
					if len(users) > 2 {
						u.Scoreboard.Set(7, text.Colourf("<grey>3. %v</grey> <black>-</black> <emerald>%v</emerald>", lo.If(users[2].Data.Cosmetics.Nickname != "", users[2].Data.Cosmetics.Nickname).Else(users[2].Data.Username), users[2].GameInfo.BuildFFA.Kills))
					}
					u.Scoreboard.Set(8, "§5")
					u.Scoreboard.Set(9, text.Colourf("<grey>Kills:</grey> <emerald>%v</emerald>", u.GameInfo.BuildFFA.Kills))
					u.Scoreboard.Set(10, text.Colourf("<grey>Position:</grey> <emerald>#%v</emerald>", slices.Index(users, u)+1))
					u.Scoreboard.Set(11, "§6")
					u.Scoreboard.Set(12, font.Transform(server.IP))
					u.SendScoreboard(5)
				}, tx)
			})
		}
	}()
}

func (b *BuildFFA) Type() game.TypeGame {
	return game.TypeBuildFFA
}

func (b *BuildFFA) Maps() []string {
	return []string{"./maps/BFFA"}
}

func (b *BuildFFA) MapConfig() game.MapData {
	return game.Maps[b.Maps()[0]]
}

func (b *BuildFFA) Handler() player.Handler {
	return Handler{}
}

func (b *BuildFFA) Reward(_ *player.Player, tx *world.Tx) (before, after int, mvp bool) {
	return 0, 0, false
}

func (b *BuildFFA) Punish(_ *player.Player, tx *world.Tx) (before, after int) {
	return 0, 0
}
