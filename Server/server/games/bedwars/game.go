package bedwars

import (
	"fmt"
	"math/rand"
	"server/server"
	"server/server/database"
	"server/server/game"
	"server/server/games/lobby"
	language2 "server/server/language"
	"server/server/user"
	"server/server/utils"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player/title"

	"github.com/df-mc/dragonfly/server/world"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/player"
)

const startingInDuration = 1 * time.Second

var Games = make(map[uuid.UUID]*BedWars)

type BedWars struct {
	*game.Game

	TeamSize    int
	TeamCount   int
	UsersToJoin []string

	typeGame   game.TypeGame
	isCustom   bool
	mapIndex   int
	startingIn time.Duration
}

func NewBedWars(typeGame game.TypeGame, teamSize int, teamCount int, isCustom bool) *BedWars {
	newId := uuid.New()
	Games[newId] = &BedWars{
		TeamSize:   teamSize,
		TeamCount:  teamCount,
		typeGame:   typeGame,
		isCustom:   isCustom,
		startingIn: startingInDuration,
	}
	g := Games[newId]

	g.mapIndex = rand.Intn(len(g.Maps()))
	mapName := g.Maps()[g.mapIndex]
	g.Game = game.NewGame(newId, utils.Panics(server.WorldManager.World(mapName)), "")
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		for range ticker.C {
			switch g.Stage() {
			case game.Waiting:
				if len(g.OriginalPlayers()) == teamSize*teamCount {
					g.SetStage(game.Starting)
				} else {
					g.ForEachActivePlayer(func(pl *player.Player) {
						sendWaitingScoreboard(pl, g)
					})
				}
				break
			case game.Starting:
				if len(g.OriginalPlayers()) != teamSize*teamCount {
					g.SetStage(game.Waiting)
					g.startingIn = startingInDuration
				} else {
					if g.startingIn == 0 {
						g.SetStage(game.Running)
						g.ForEachActivePlayer(func(pl *player.Player) {
							team := g.PlayerTeam(pl)
							pl.Teleport(g.MapConfig().TeamSpawnPoints[team.ID()])
							giveKit(pl, g)

							u := user.LookupPlayer(pl)
							if g.typeGame == game.TypeBedWars {
								u.Data.Games.BedWars.GamesPlayed++
							} else {
								u.Data.Games.BedFight.GamesPlayed++
							}
						})
					} else {
						g.ForEachActivePlayer(func(pl *player.Player) {
							sendStartingScoreboard(pl, g)
						})
						g.startingIn -= 250 * time.Millisecond
					}
				}
				break
			case game.Running:
				if g.WinningTeam() != nil {
					g.SetStage(game.Ending)

					g.ForEachActivePlayer(func(pl *player.Player) {
						if g.WinningTeam().Contains(pl) {
							g.Reward(pl)
							pl.SendTitle(title.New(text.Colourf(language2.Translate(pl).BedWars.YouWonTitle)).WithSubtitle(text.Colourf(language2.Translate(pl).BedWars.TeamWonSubTitle, g.WinningTeam().Color(), strings.ToUpper(g.WinningTeam().Color()), g.WinningTeam().Color())))
						} else {
							g.Punish(pl)
							pl.SendTitle(title.New(text.Colourf(language2.Translate(pl).BedWars.YouLostTitle)).WithSubtitle(text.Colourf(language2.Translate(pl).BedWars.TeamWonSubTitle, g.WinningTeam().Color(), strings.ToUpper(g.WinningTeam().Color()), g.WinningTeam().Color())))
						}
					})
				} else {
					g.ForEachActivePlayer(func(pl *player.Player) {
						sendRunningScoreboard(pl, g)
					})
				}
				break
			case game.Ending:
				time.AfterFunc(5*time.Second, func() {
					g.SetStage(game.Terminated)
				})
				break
			case game.Terminated:
				ticker.Stop()

				g.ForEachOriginalPlayer(func(pl *player.Player) {
					u := user.LookupPlayer(pl)
					u.GameInfo = user.GameRuntimeData{}
				})

				g.World().Exec(func(tx *world.Tx) {
					for e := range tx.Players() {
						pl := e.(*player.Player)
						pl.Handler().HandleQuit(pl)
						tx.RemoveEntity(pl)
						server.MCServer.World().Exec(func(tx *world.Tx) {
							tx.AddEntity(pl.H())
						})

						lobby.Join(pl)
					}
				})
				break
			default:
				panic("unknown stage")
			}
		}
	}()

	return g
}

func (b *BedWars) Type() game.TypeGame {
	return b.typeGame
}

func (b *BedWars) Maps() []string {
	if b.typeGame == game.TypeBedWars {
		return []string{"Invasion"}
	}
	return []string{"BedFight1"}
}

func (b *BedWars) MapConfig() game.MapData {
	return game.Maps[fmt.Sprintf("./maps/%v", b.Maps()[b.mapIndex])]
}

func (b *BedWars) Handler() player.Handler {
	return Handler{}
}

func (b *BedWars) Reward(pl *player.Player) {
	var sorted []*player.Player
	b.ForEachOriginalPlayer(func(p *player.Player) {
		sorted = append(sorted, pl)
	})

	slices.SortFunc(sorted, func(a, b *player.Player) int {
		ua := user.LookupPlayer(a)
		ub := user.LookupPlayer(b)

		return ub.GameInfo.BedWars.Kills - ua.GameInfo.BedWars.FinalKills
	})

	mostKills := user.LookupPlayer(sorted[0]).GameInfo.BedWars.Kills + user.LookupPlayer(sorted[0]).GameInfo.BedWars.FinalKills

	u := user.LookupPlayer(pl)
	if u.GameInfo.BedWars.Kills+u.GameInfo.BedWars.FinalKills == mostKills {
		switch u.Data.Statistics.ELORank() {
		case database.Bronze:
			u.Data.Statistics.ELO += 25
			break
		case database.Silver, database.Gold, database.Platinum, database.Diamond:
			u.Data.Statistics.ELO += 20
			break
		case database.Emerald, database.Sapphire, database.Ruby, database.Crystal:
			u.Data.Statistics.ELO += 15
			break
		case database.Opal, database.Amethyst, database.Obsidian, database.Aventurine:
			u.Data.Statistics.ELO += 10
			break
		case database.Quartz, database.Topaz, database.DarkMatter:
			u.Data.Statistics.ELO += 5
			break
		}
	}

	switch u.Data.Statistics.ELORank() {
	case database.Bronze:
		u.Data.Statistics.ELO += 45
		break
	case database.Silver, database.Gold:
		u.Data.Statistics.ELO += 40
		break
	case database.Platinum, database.Diamond:
		u.Data.Statistics.ELO += 35
		break
	case database.Emerald, database.Sapphire:
		u.Data.Statistics.ELO += 30
		break
	case database.Ruby, database.Crystal:
		u.Data.Statistics.ELO += 25
		break
	case database.Opal, database.Amethyst:
		u.Data.Statistics.ELO += 20
		break
	case database.Obsidian, database.Aventurine:
		u.Data.Statistics.ELO += 15
		break
	case database.Quartz, database.Topaz:
		u.Data.Statistics.ELO += 10
		break
	case database.DarkMatter:
		u.Data.Statistics.ELO += 5
		break
	}

	if b.typeGame == game.TypeBedWars {
		u.Data.Games.BedWars.Wins++
		u.Data.Games.BedWars.WinStreak++
	} else {
		u.Data.Games.BedFight.Wins++
		u.Data.Games.BedFight.WinStreak++
	}
}

func (b *BedWars) Punish(pl *player.Player) {
	u := user.LookupPlayer(pl)
	switch u.Data.Statistics.ELORank() {
	case database.Bronze, database.Silver:
		u.Data.Statistics.ELO -= 10
		break
	case database.Gold, database.Platinum:
		u.Data.Statistics.ELO -= 15
		break
	case database.Diamond, database.Emerald:
		u.Data.Statistics.ELO -= 20
		break
	case database.Sapphire, database.Ruby:
		u.Data.Statistics.ELO -= 25
		break
	case database.Crystal, database.Opal:
		u.Data.Statistics.ELO -= 30
		break
	case database.Amethyst, database.Obsidian:
		u.Data.Statistics.ELO -= 35
		break
	case database.Aventurine, database.Quartz:
		u.Data.Statistics.ELO -= 40
		break
	case database.Topaz, database.DarkMatter:
		u.Data.Statistics.ELO -= 45
		break
	}

	if b.typeGame == game.TypeBedWars {
		u.Data.Games.BedWars.Losses++
		u.Data.Games.BedWars.WinStreak = 0
	} else {
		u.Data.Games.BedFight.Losses++
		u.Data.Games.BedFight.WinStreak = 0
	}
}

func sendWaitingScoreboard(pl *player.Player, g *BedWars) {
	u := user.LookupPlayer(pl)
	u.Scoreboard.Set(0, text.Colourf("          <yellow>▷ <white>Season 1</white> ◁</yellow>"))
	u.Scoreboard.Set(1, "§0")
	u.Scoreboard.Set(2, text.Colourf("<white>Map:</white> <green>%v</green>", g.MapConfig().Name))
	u.Scoreboard.Set(3, text.Colourf("<white>Players:</white> <green>%v/%v</green>", len(g.ActivePlayers()), g.TeamCount*g.TeamSize))
	u.Scoreboard.Set(4, "§1")
	u.Scoreboard.Set(5, text.Colourf("<white>Phase:</white> <green>%v</green>", "Waiting for players..."))
	u.Scoreboard.Set(6, "§2")
	i := 7
	if g.typeGame == game.TypeBedWars {
		u.Scoreboard.Set(i, text.Colourf("<white>Mode:</white> <green>%v</green>", "Ranked"))
		i++
		u.Scoreboard.Set(i, "§3")
		i++
	}
	u.Scoreboard.Set(i, text.Colourf("<yellow>ELIAGIC.CLUB</yellow>"))
	pl.SendScoreboard(u.Scoreboard)
}

func sendStartingScoreboard(pl *player.Player, g *BedWars) {
	u := user.LookupPlayer(pl)
	u.Scoreboard.Set(0, text.Colourf("          <yellow>▷ <white>Season 1</white> ◁</yellow>"))
	u.Scoreboard.Set(1, "§0")
	u.Scoreboard.Set(2, text.Colourf("<white>Map:</white> <green>%v</green>", g.MapConfig().Name))
	u.Scoreboard.Set(3, text.Colourf("<white>Players:</white> <green>%v/%v</green>", len(g.ActivePlayers()), g.TeamCount*g.TeamSize))
	u.Scoreboard.Set(4, "§1")
	u.Scoreboard.Set(5, text.Colourf("<grey>Phase:</grey> <emerald>Starting in <yellow>%.2f</yellow> seconds</emerald>", g.startingIn.Seconds()))
	u.Scoreboard.Set(6, "§2")
	i := 7
	if g.typeGame == game.TypeBedWars {
		if g.isCustom {
			u.Scoreboard.Set(i, text.Colourf("<white>Mode:</white> <green>%v</green>", "Custom"))
		} else {
			u.Scoreboard.Set(i, text.Colourf("<white>Mode:</white> <green>%v</green>", "Ranked"))
		}
		i++
		u.Scoreboard.Set(i, "§3")
		i++
	}
	u.Scoreboard.Set(i, text.Colourf("<yellow>ELIAGIC.CLUB</yellow>"))
	pl.SendScoreboard(u.Scoreboard)

}

func sendRunningScoreboard(pl *player.Player, g *BedWars) {
	u := user.LookupPlayer(pl)
	u.Scoreboard.Set(0, text.Colourf("          <yellow>▷ <white>Season 1</white> ◁</yellow>"))
	u.Scoreboard.Set(1, "§0")
	i := 2
	for _, t := range g.Teams() {
		var statusStr string
		switch t.Status {
		case game.BedExists:
			statusStr = "<green>✔</green>"
			break
		case game.BedBroken:
			statusStr = strconv.Itoa(t.CountActivePlayers())
			break
		case game.TeamDead:
			statusStr = "<red>✖</red>"
			break
		}

		u.Scoreboard.Set(i, text.Colourf(
			"<%v>%v</%v> <white>%v:</white> <bold>%v</bold>",
			t.Color(),
			strings.ToUpper(string([]rune(t.Color())[0])),
			t.Color(),
			cases.Title(language.English).String(strings.Replace(t.Color(), "-", " ", 1)),
			statusStr,
		))

		i++
	}
	u.Scoreboard.Set(i, "§1")
	i++
	u.Scoreboard.Set(i, text.Colourf("<white>Kills:</white> <green>%v</green>", u.GameInfo.BedWars.Kills))
	i++
	u.Scoreboard.Set(i, "§2")
	i++
	u.Scoreboard.Set(i, text.Colourf("<yellow>ELIAGIC.CLUB</yellow>"))

	pl.SendScoreboard(u.Scoreboard)
}
