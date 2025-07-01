package bedwars

import (
	"fmt"
	"math"
	"math/rand"
	"server/server"
	"server/server/database"
	"server/server/font"
	"server/server/game"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"

	"github.com/df-mc/dragonfly/server/entity"

	"github.com/samber/lo"
	language2 "golang.org/x/text/language"

	"github.com/df-mc/dragonfly/server/item/enchantment"

	"github.com/df-mc/dragonfly/server/item"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server/player/title"

	"github.com/df-mc/dragonfly/server/world"

	"github.com/google/uuid"
	"golang.org/x/text/cases"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/player"
)

const startingInDuration = 30 * time.Second

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

	ironGeneratorSettings    *GeneratorSettings
	goldGeneratorSettings    *GeneratorSettings
	diamondGeneratorSettings *GeneratorSettings
	emeraldGeneratorSettings *GeneratorSettings
	generators               []*GeneratorBlockType

	pickaxeTierPlayers map[uuid.UUID]int
	axeTierPlayers     map[uuid.UUID]int

	trapIgnore map[uuid.UUID]bool
}

func NewBedWars(typeGame game.TypeGame, teamSize int, teamCount int, isCustom bool) *BedWars {
	newId := uuid.New()
	Games[newId] = &BedWars{
		TeamSize:           teamSize,
		TeamCount:          teamCount,
		typeGame:           typeGame,
		isCustom:           isCustom,
		startingIn:         startingInDuration,
		pickaxeTierPlayers: make(map[uuid.UUID]int),
		axeTierPlayers:     make(map[uuid.UUID]int),
		trapIgnore:         make(map[uuid.UUID]bool),
	}
	g := Games[newId]

	g.mapIndex = rand.Intn(len(g.Maps()))
	mapName := g.Maps()[g.mapIndex]
	gameWorld := utils.Panics(server.WorldManager.World(mapName))
	gameWorld.Handle(WorldHandler{game: g})
	g.Game = game.NewGame(newId, gameWorld, "")

	go func() {
		stages := []*stage{
			{action: "<diamond>Diamond Generators</diamond>", tier: 2, dur: 6 * time.Minute},
			{action: "<emerald>Emerald Generators</emerald>", tier: 2, dur: 6 * time.Minute},
			{action: "<diamond>Diamond Generators</diamond>", tier: 3, dur: 6 * time.Minute},
			{action: "<emerald>Emerald Generators</emerald>", tier: 3, dur: 6 * time.Minute},
			{action: "<red>Bed Gone</red>", dur: 6 * time.Minute},
			{action: "<red>Sudden Death | Phase I</Red>", tier: 1, dur: 10 * time.Minute},
			{action: "<red>Sudden Death | Phase II</Red>", tier: 2, dur: 3 * time.Minute},
			{action: "<red>Sudden Death | Phase III</Red>", tier: 3, dur: 3 * time.Minute},
			{action: "<red>Game Ends</red>", dur: 3 * time.Minute},
		}
		ticker := time.NewTicker(100 * time.Millisecond)
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
			case game.Starting:
				if len(g.OriginalPlayers()) != teamSize*teamCount {
					g.startingIn = startingInDuration
					g.SetStage(game.Waiting)
				} else {
					cds := map[int]string{
						30: "<green>30</green>",
						20: "<gold>20</gold>",
						10: "<gold>10</gold>",
						5:  "<red>5</red>",
						4:  "<red>4</red>",
						3:  "<red>3</red>",
						2:  "<red>2</red>",
						1:  "<red>1</red>",
					}

					if cd := cds[int(g.startingIn.Seconds())]; cd != "" && math.Mod(g.startingIn.Seconds(), 1) == 0 {
						g.ForEachOriginalPlayer(func(pl *player.Player) {
							pl.Message(text.Colourf(language.Translate(pl).Game.Countdown, cd))
							if int(g.startingIn.Seconds()) <= 5 {
								pl.SendTitle(title.New(text.Colourf(cd)))
							}
							pl.PlaySound(sound.Click{})
						})
					}

					if g.startingIn == 0 {
						g.World().Exec(func(tx *world.Tx) {
							g.initBedWarsFeatures(tx)
						})

						if len(g.UsersToJoin) != 0 {
							g.ForEachActivePlayer(func(pl *player.Player) {
								g.RemovePlayerFromTeam(pl)
							})

							for _, userId := range g.UsersToJoin {
								u := user.GetUserByUserID(userId)
								g.AddPlayerToTeam(u.Player(), g.TeamSize)
							}
						}

						g.ForEachActivePlayer(func(pl *player.Player) {
							u := user.GetUser(pl)
							team := g.PlayerTeam(pl)

							pl.SetNameTag(database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Colour()).Name(u.Data))
							pl.Teleport(g.MapConfig().TeamSpawnPoints[team.ID()])
							giveKit(pl, g)

							if g.typeGame == game.TypeBedWars {
								u.Data.Games.BedWars.GamesPlayed++
							} else {
								u.Data.Games.BedFight.GamesPlayed++
							}

							for _, gen := range g.generators {
								gen.Active = true
							}

							pl.SendTitle(title.New(text.Colourf("<green>GO!</green>")))
							if g.typeGame == game.TypeBedWars {
								pl.Message(text.Colourf(language.Translate(pl).BedWars.TutorialMessage))
							}
						})

						g.SetStage(game.Running)
					} else {
						g.ForEachActivePlayer(func(pl *player.Player) {
							sendStartingScoreboard(pl, g)
						})
						g.startingIn -= 100 * time.Millisecond
					}
				}
			case game.Running:
				suddenDeathTicker := time.NewTicker(700 * time.Millisecond)
				if g.WinningTeam() != nil {
					suddenDeathTicker.Stop()
					g.ForEachActivePlayer(func(pl *player.Player) {
						if g.WinningTeam().Contains(pl) {
							go g.Reward(pl)
							pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.VictoryTitle)))
						} else {
							g.Punish(pl)
							pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.DefeatTitle)))
						}

						var name string
						var mostKills int
						g.WinningTeam().ForEachPlayer(pl.Tx(), func(p *player.Player) {
							uwt := user.GetUser(p)
							if name == "" || mostKills < uwt.GameInfo.TotalBWKills() {
								name = database.LobbyNameDisplay.Name(uwt.Data)
								mostKills = uwt.GameInfo.TotalBWKills()
							}
						})

						var sorted []*player.Player
						for _, e := range g.OriginalPlayers() {
							if p, ok := e.Entity(pl.Tx()); ok {
								sorted = append(sorted, p.(*player.Player))
							}
						}

						slices.SortFunc(sorted, func(a, b *player.Player) int {
							ua := user.GetUser(a)
							ub := user.GetUser(b)

							return ub.GameInfo.TotalBWKills() - ua.GameInfo.TotalBWKills()
						})

						var l1, l2, l3, l4 string

						if name != "" {
							l1 = text.Colourf("<%v>%v</%v> <grey>-</grey> %v", g.WinningTeam().Colour(), strings.ToUpper(g.WinningTeam().Colour()), g.WinningTeam().Colour(), name)
						}

						if len(sorted) > 0 {
							u := user.GetUser(sorted[0])
							l2 = text.Colourf("<yellow>1st Killer</yellow> <grey>-</grey> %v <grey>- %v</grey>", database.LobbyNameDisplay.Name(u.Data), u.GameInfo.TotalBWKills())
						}
						if len(sorted) > 1 {
							u := user.GetUser(sorted[1])
							l3 = text.Colourf("<gold>2nd Killer</gold> <grey>-</grey> %v <grey>- %v</grey>", database.LobbyNameDisplay.Name(u.Data), u.GameInfo.TotalBWKills())
						}
						if len(sorted) > 2 {
							u := user.GetUser(sorted[2])
							l4 = text.Colourf("<red>3rd Killer</red> <grey>-</grey> %v <grey>- %v</grey>", database.LobbyNameDisplay.Name(u.Data), u.GameInfo.TotalBWKills())
						}

						pl.Message(text.Colourf(
							`<green>============================================================</green>
                                    <bold>Bed Wars</bold>

%v%v

%v%v
%v%v
%v%v

<green>============================================================</green>`,
							strings.Repeat(" ", 80-len(l1)),
							l1,
							strings.Repeat(" ", 90-len(l2)),
							l2,
							strings.Repeat(" ", 90-len(l3)),
							l3,
							strings.Repeat(" ", 90-len(l4)),
							l4,
						))
					})

					g.SetStage(game.Ending)
				} else {
					currentStage := stages[0]
					g.ForEachActivePlayer(func(pl *player.Player) {
						sendRunningScoreboard(pl, g, currentStage)
					})

					currentStage.dur -= 100 * time.Millisecond
					if currentStage.dur == 0 {
						switch currentStage.action {
						case "<diamond>Diamond Generators</diamond>", "<emerald>Emerald Generators</emerald>":
							if currentStage.action == "<diamond>Diamond Generators</diamond>" {
								g.diamondGeneratorSettings.Tier++
								g.diamondGeneratorSettings.SpawnRate = lo.If(currentStage.tier == 2, 23*time.Second).Else(12 * time.Second)
							} else {
								g.emeraldGeneratorSettings.Tier++
								g.emeraldGeneratorSettings.SpawnRate = lo.If(currentStage.tier == 2, 40*time.Second).Else(27 * time.Second)
							}
							g.World().Exec(func(tx *world.Tx) {
								for e := range tx.Players() {
									pl := e.(*player.Player)
									pl.Message(text.Colourf(language.Translate(pl).BedWars.GeneratorUpgraded, currentStage.action, currentStage.tier))
								}
							})
						case "<red>Bed Gone</red>":
							g.World().Exec(func(tx *world.Tx) {
								for e := range tx.Players() {
									pl := e.(*player.Player)
									pl.Message(text.Colourf(language.Translate(pl).BedWars.BedGone))
									pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.BedBreakTitle)).WithSubtitle(text.Colourf(language.Translate(pl).BedWars.BedBreakSubTitle)))
								}
								for _, team := range g.Teams() {
									tx.SetBlock(cube.PosFromVec3(g.MapConfig().BedPositions[team.ID()*2]), block.Air{}, nil)
									tx.SetBlock(cube.PosFromVec3(g.MapConfig().BedPositions[team.ID()*2+1]), block.Air{}, nil)

									g.Teams()[team.ID()].Status = game.BedBroken
								}

								g.playBedBrokenSound(tx)
							})
						case "<red>Sudden Death | Phase I</Red>", "<red>Sudden Death | Phase II</Red>", "<red>Sudden Death | Phase III</Red>":
							g.World().Exec(func(tx *world.Tx) {
								for e := range tx.Players() {
									pl := e.(*player.Player)
									pl.Message(text.Colourf(language.Translate(pl).BedWars.SuddenDeath))
									pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.SuddenDeathTitle)))
								}
							})

							go func() {
								for range suddenDeathTicker.C {
									g.ForEachActivePlayer(func(pl *player.Player) {
										if utils.RandChance(20 * currentStage.tier) {
											ePos := pl.Position().Add(mgl64.Vec3{
												rand.Float64()*15 - 7.5,
												30,
												rand.Float64()*15 - 7.5,
											})

											if utils.RandChance(50) {
												pl.Tx().AddEntity(entity.NewTNT(world.EntitySpawnOpts{Position: ePos}, 10*time.Second))
											} else {
												NewFireball(ePos, pl.Tx())
											}
										}
									})
								}
							}()
						case "<red>Game Ends</red>":
							suddenDeathTicker.Stop()
							g.World().Exec(func(tx *world.Tx) {
								for e := range tx.Players() {
									pl := e.(*player.Player)
									pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.DrawTitle)))
								}
							})
							g.SetStage(game.Ending)
						}
						stages = stages[1:]
					}
				}
			case game.Ending:
				time.AfterFunc(5*time.Second, func() {
					g.SetStage(game.Terminated)
				})
			case game.Terminated:
				ticker.Stop()

				g.ForEachOriginalPlayer(func(pl *player.Player) {
					u := user.GetUser(pl)
					u.GameInfo = user.GameRuntimeData{}
				})

				<-g.World().Exec(func(tx *world.Tx) {
					for e := range tx.Players() {
						pl := e.(*player.Player)
						pl.Handler().HandleQuit(pl)
						tx.RemoveEntity(pl)
						server.MCServer.World().Exec(func(tx *world.Tx) {
							tx.AddEntity(pl.H())
						})
					}
				})

				utils.Panic(g.World().Close())
				delete(Games, g.ID())
			default:
				panic("unknown stage")
			}
		}
	}()

	return g
}

func (b *BedWars) initBedWarsFeatures(tx *world.Tx) {
	for _, pos := range b.MapConfig().ShopVillagerPositions {
		t := b.NearestTeam(pos)
		v := NewItemsVillager(pos, text.Colourf("<green>Shop Villager</green>"), b, t, tx)
		v.LookAt(b.MapConfig().TeamSpawnPoints[t.ID()].Sub(mgl64.Vec3{0, 2, 0}), tx)
	}
	for _, pos := range b.MapConfig().UpgradesVillagerPositions {
		t := b.NearestTeam(pos)
		v := NewUpgradesVillager(pos, text.Colourf("<green>Upgrades Villager</green>"), b, t, tx)
		v.LookAt(b.MapConfig().TeamSpawnPoints[t.ID()].Sub(mgl64.Vec3{0, 2, 0}), tx)
	}

	b.ironGeneratorSettings = &GeneratorSettings{
		Game:      b,
		Resource:  Iron,
		Tier:      1,
		Cap:       48,
		SpawnRate: 400 * time.Millisecond,
	}
	b.goldGeneratorSettings = &GeneratorSettings{
		Game:      b,
		Resource:  Gold,
		Tier:      1,
		Cap:       12,
		SpawnRate: 4 * time.Second,
	}
	b.diamondGeneratorSettings = &GeneratorSettings{
		Resource:  Diamond,
		Tier:      1,
		Name:      text.Colourf("<bold><diamond>Diamond</diamond></bold>"),
		Cap:       8,
		SpawnRate: 30 * time.Second,
	}
	b.emeraldGeneratorSettings = &GeneratorSettings{
		Resource:  Emerald,
		Tier:      1,
		Name:      text.Colourf("<bold><emerald>Emerald</emerald></bold>"),
		Cap:       6,
		SpawnRate: 66 * time.Second,
	}

	for _, pos := range b.MapConfig().IronGenerators {
		b.generators = append(b.generators, b.ironGeneratorSettings.New(pos, tx), b.goldGeneratorSettings.New(pos, tx))
	}

	for _, pos := range b.MapConfig().DiamondGenerators {
		b.generators = append(b.generators, b.diamondGeneratorSettings.New(pos, tx))
	}

	for _, pos := range b.MapConfig().EmeraldGenerators {
		b.generators = append(b.generators, b.emeraldGeneratorSettings.New(pos, tx))
	}

	for _, pos := range b.MapConfig().ChestPositions {
		tx.AddEntity(entity.NewText(text.Colourf("<yellow>PUNCH TO</yellow>\n<yellow>DEPOSIT</yellow>"), pos))
	}

	for _, pos := range b.MapConfig().EnderChestPositions {
		tx.AddEntity(entity.NewText(text.Colourf("<yellow>PUNCH TO</yellow>\n<yellow>DEPOSIT</yellow>"), pos))
	}
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
	return PlayerHandler{}
}

func (b *BedWars) Reward(pl *player.Player) {
	var sorted []*player.Player
	b.ForEachOriginalPlayer(func(p *player.Player) {
		sorted = append(sorted, pl)
	})

	slices.SortFunc(sorted, func(a, b *player.Player) int {
		ua := user.GetUser(a)
		ub := user.GetUser(b)

		return ub.GameInfo.BedWars.Kills - ua.GameInfo.BedWars.Kills
	})

	mostKills := user.GetUser(sorted[0]).GameInfo.TotalBWKills()

	u := user.GetUser(pl)
	if u.GameInfo.TotalBWKills() == mostKills {
		u.Data.Games.BedWars.MVPCount++
		switch u.Data.Statistics.ELORank() {
		case database.Bronze:
			u.Data.Statistics.ELO += 25
		case database.Silver, database.Gold, database.Platinum, database.Diamond:
			u.Data.Statistics.ELO += 20
		case database.Emerald, database.Sapphire, database.Ruby, database.Crystal:
			u.Data.Statistics.ELO += 15
		case database.Opal, database.Amethyst, database.Obsidian, database.Aventurine:
			u.Data.Statistics.ELO += 10
		case database.Quartz, database.Topaz, database.DarkMatter:
			u.Data.Statistics.ELO += 5
		}
	}

	switch u.Data.Statistics.ELORank() {
	case database.Bronze:
		u.Data.Statistics.ELO += 45
	case database.Silver, database.Gold:
		u.Data.Statistics.ELO += 40
	case database.Platinum, database.Diamond:
		u.Data.Statistics.ELO += 35
	case database.Emerald, database.Sapphire:
		u.Data.Statistics.ELO += 30
	case database.Ruby, database.Crystal:
		u.Data.Statistics.ELO += 25
	case database.Opal, database.Amethyst:
		u.Data.Statistics.ELO += 20
	case database.Obsidian, database.Aventurine:
		u.Data.Statistics.ELO += 15
	case database.Quartz, database.Topaz:
		u.Data.Statistics.ELO += 10
	case database.DarkMatter:
		u.Data.Statistics.ELO += 5
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
	u := user.GetUser(pl)
	switch u.Data.Statistics.ELORank() {
	case database.Bronze, database.Silver:
		u.Data.Statistics.ELO -= 10
	case database.Gold, database.Platinum:
		u.Data.Statistics.ELO -= 15
	case database.Diamond, database.Emerald:
		u.Data.Statistics.ELO -= 20
	case database.Sapphire, database.Ruby:
		u.Data.Statistics.ELO -= 25
	case database.Crystal, database.Opal:
		u.Data.Statistics.ELO -= 30
	case database.Amethyst, database.Obsidian:
		u.Data.Statistics.ELO -= 35
	case database.Aventurine, database.Quartz:
		u.Data.Statistics.ELO -= 40
	case database.Topaz, database.DarkMatter:
		u.Data.Statistics.ELO -= 45
	}

	if b.typeGame == game.TypeBedWars {
		u.Data.Games.BedWars.Losses++
		u.Data.Games.BedWars.WinStreak = 0
	} else {
		u.Data.Games.BedFight.Losses++
		u.Data.Games.BedFight.WinStreak = 0
	}
}

func (b *BedWars) NearestGenerator(pos mgl64.Vec3, resource Resource) *GeneratorBlockType {
	var nearestGen *GeneratorBlockType
	for _, gen := range b.generators {
		if gen.Resource == resource && (nearestGen == nil || utils.Distance(pos, gen.Position()) < utils.Distance(pos, nearestGen.Position())) {
			nearestGen = gen
		}
	}
	return nearestGen
}

func (b *BedWars) NearestTeam(pos mgl64.Vec3) *game.Team {
	var nearestTeam *game.Team
	for _, t := range b.Teams() {
		if nearestTeam == nil || utils.Distance(pos, b.MapConfig().TeamSpawnPoints[t.ID()]) < utils.Distance(pos, b.MapConfig().TeamSpawnPoints[nearestTeam.ID()]) {
			nearestTeam = t
		}
	}
	return nearestTeam
}

func (b *BedWars) NearestEnemyTeam(team *game.Team, pos mgl64.Vec3) *game.Team {
	var nearestTeam *game.Team
	for _, t := range b.Teams() {
		if t != team && (nearestTeam == nil || utils.Distance(pos, b.MapConfig().TeamSpawnPoints[t.ID()]) < utils.Distance(pos, b.MapConfig().TeamSpawnPoints[nearestTeam.ID()])) {
			nearestTeam = t
		}
	}
	return nearestTeam
}

func (b *BedWars) buyItem(pl *player.Player, s item.Stack) bool {
	u := user.GetUser(pl)
	if canAfford(pl, s) {
		resource, cost := getCost(s)
		_ = pl.Inventory().RemoveItem(item.NewStack(resource.Item(), cost))

		addItem := func() bool {
			if n, err := u.AddItemWithHBConfig(-1, s.WithLore()); err != nil {
				_ = pl.Inventory().RemoveItem(item.NewStack(s.Item(), n))
				_, _ = pl.Inventory().AddItem(item.NewStack(resource.Item(), cost))
				return false
			}

			return true
		}

		if boots, ok := s.Item().(item.Boots); ok {
			t := b.PlayerTeam(pl)
			pl.Armour().Set(
				item.NewStack(item.Helmet{Tier: boots.Tier}, 1).AsUnbreakable(),
				item.NewStack(item.Chestplate{Tier: boots.Tier}, 1).AsUnbreakable(),
				item.NewStack(item.Leggings{Tier: boots.Tier}, 1).AsUnbreakable(),
				item.NewStack(item.Boots{Tier: boots.Tier}, 1).AsUnbreakable(),
			)

			if t.Upgrades.Protection != 0 {
				for slot, stack := range pl.Armour().Items() {
					utils.Panic(pl.Armour().Inventory().SetItem(slot, stack.WithEnchantments(item.NewEnchantment(enchantment.Protection, t.Upgrades.Protection))))
				}
			}
		} else if _, ok := s.Item().(item.Pickaxe); ok {
			b.pickaxeTierPlayers[pl.UUID()]++

			var flag bool
			for slot, stack := range pl.Inventory().Items() {
				if _, ok := stack.Item().(item.Pickaxe); ok {
					utils.Panic(pl.Inventory().SetItem(slot, s))
					flag = true
				}
			}

			if !flag {
				return addItem()
			}
		} else if _, ok := s.Item().(item.Axe); ok {
			b.axeTierPlayers[pl.UUID()]++

			var flag bool
			for slot, stack := range pl.Inventory().Items() {
				if _, ok := stack.Item().(item.Axe); ok {
					utils.Panic(pl.Inventory().SetItem(slot, s))
					flag = true
				}
			}

			if !flag {
				return addItem()
			}
		} else {
			return addItem()
		}

		return true
	}
	return false
}

func (b *BedWars) buyUpgrade(pl *player.Player, s item.Stack) bool {
	if canAfford(pl, s) {
		resource, cost := getCost(s)
		_ = pl.Inventory().RemoveItem(item.NewStack(resource.Item(), cost))
		return true
	}
	return false
}

func (b *BedWars) playBedBrokenSound(tx *world.Tx) {
	for e := range tx.Players() {
		pl := e.(*player.Player)
		user.GetUser(pl).PlaySound("mob.enderdragon.growl", 1, 1)
	}
}

func sendWaitingScoreboard(pl *player.Player, g *BedWars) {
	u := user.GetUser(pl)
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
	u.Scoreboard.Set(i, font.Transform(server.IP))
	u.SendScoreboard(12)
}

func sendStartingScoreboard(pl *player.Player, g *BedWars) {
	u := user.GetUser(pl)
	u.Scoreboard.Set(1, "§0")
	u.Scoreboard.Set(2, text.Colourf("<white>Map:</white> <green>%v</green>", g.MapConfig().Name))
	u.Scoreboard.Set(3, text.Colourf("<white>Players:</white> <green>%v/%v</green>", len(g.ActivePlayers()), g.TeamCount*g.TeamSize))
	u.Scoreboard.Set(4, "§1")
	u.Scoreboard.Set(5, text.Colourf("<grey>Phase:</grey> <emerald>Starting in <yellow>%.1f</yellow> seconds</emerald>", g.startingIn.Seconds()))
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
	u.Scoreboard.Set(i, font.Transform(server.IP))
	u.SendScoreboard(14)
}

func sendRunningScoreboard(pl *player.Player, g *BedWars, stage *stage) {
	u := user.GetUser(pl)
	u.Scoreboard.Set(1, "§0")
	i := 2
	if g.Type() == game.TypeBedWars {
		u.Scoreboard.Set(i, text.Colourf("%v in <green>%d:%02d</green>", strings.Replace(stage.action, "Generators", utils.IntToRoman(stage.tier), 1), int(stage.dur.Seconds())/60, int(stage.dur.Seconds())%60))
		i++
		u.Scoreboard.Set(i, "§1")
		i++
	}
	for _, t := range g.Teams() {
		var statusStr string
		switch t.Status {
		case game.BedExists:
			statusStr = "\uE100"
		case game.BedBroken:
			statusStr = strconv.Itoa(t.CountActivePlayers())
		case game.TeamDead:
			statusStr = "\uE101"
		}

		u.Scoreboard.Set(i, text.Colourf(
			"<%v>%v</%v> <white>%v:</white> %v",
			t.Colour(),
			strings.ToUpper(string([]rune(t.Colour())[0])),
			t.Colour(),
			cases.Title(language2.English).String(strings.Replace(t.Colour(), "-", " ", 1)),
			statusStr,
		))

		i++
	}
	u.Scoreboard.Set(i, "§2")
	i++
	u.Scoreboard.Set(i, text.Colourf("<white>Kills:</white> <green>%v</green>", u.GameInfo.BedWars.Kills))
	i++
	if g.Type() == game.TypeBedWars {
		u.Scoreboard.Set(i, text.Colourf("<white>Final Kills:</white> <green>%v</green>", u.GameInfo.BedWars.FinalKills))
		i++
		u.Scoreboard.Set(i, text.Colourf("<white>Beds Broken:</white> <green>%v</green>", u.GameInfo.BedWars.BedsBroken))
		i++
	}
	u.Scoreboard.Set(i, "§3")
	i++
	u.Scoreboard.Set(i, font.Transform(server.IP))

	u.SendScoreboard(6)
}

type stage struct {
	action string
	tier   int
	dur    time.Duration
}
