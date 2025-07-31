package bedwars

import (
	"fmt"
	core "server/server"
	"server/server/blocks/bed"
	"server/server/database"
	"server/server/game"
	"server/server/inv"
	"server/server/language"
	"server/server/listener"
	"server/server/living"
	"server/server/user"
	"server/server/utils"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/df-mc/dragonfly/server/event"

	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/item/inventory"

	"github.com/df-mc/dragonfly/server/item/enchantment"

	"github.com/df-mc/dragonfly/server/item/potion"

	"github.com/df-mc/dragonfly/server/player/title"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var blocksPlaced = make(map[string]world.Block)

type PlayerHandler struct {
	player.NopHandler

	game *BedWars
}

func Join(pl *player.Player, tx *world.Tx, teamSize int, teamCount int, typeGame game.TypeGame, isCustom bool, bwGame *BedWars) {
	if bwGame == nil {
		for _, g := range Games {
			if g.Type() == typeGame && g.Stage() == game.Waiting {
				bwGame = g
				break
			}
		}
	}

	if bwGame == nil {
		bwGame = NewBedWars(typeGame, teamSize, teamCount, isCustom)
	}

	core.Players[pl.UUID()] = pl.Name()

	tx.RemoveEntity(pl)

	bwGame.World().Exec(func(tx *world.Tx) {
		tx.AddEntity(pl.H())
	})

	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()

	u := user.GetUser(pl)
	u.Game = bwGame.Game
	switch typeGame {
	case game.TypeBedWars:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDWARS</yellow></bold>"))
	case game.TypeBedFight:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDFIGHT</yellow></bold>"))
	default:
		panic("Unhandled game type")
	}

	bwGame.AddPlayerToTeam(pl, teamSize, typeGame)

	pl.Teleport(bwGame.MapConfig().SpawnPoint)

	pl.Handle(PlayerHandler{game: bwGame})
	updateMenu := func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
		if shop := activeItemShops[pl.UUID()]; shop != nil {
			openedPos := utils.FetchPrivateField[atomic.Pointer[cube.Pos]](utils.Session(pl), "openedPos")
			if _, ok := ctx.Val().(*player.Player).Tx().Block(*openedPos.Load()).(block.Air); ok {
				isTools := shop.inv.ContainsItemFunc(1, func(stack item.Stack) bool {
					_, ok := stack.Item().(item.Shears)
					return ok
				})
				it43, _ := shop.inv.Item(43)
				isQuickBuy := !it43.Empty()

				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					if isTools {
						shop.inv.Clear()
						for i, it := range shop.Tools() {
							_ = shop.inv.SetItem(i, it)
						}
					} else if isQuickBuy {
						shop.inv.Clear()
						for i, it := range shop.itemShopDashboard(true) {
							_ = shop.inv.SetItem(i, it)
						}
					}
				}()
			}
		}
	}
	pl.Inventory().Handle(inv.ChestUIHandler{Inventory: pl.Inventory(), Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		updateMenu,
		updateMenu,
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			updateMenu(ctx, slot, stack, inv)
			_, ok := stack.Item().(item.Tool)
			if ok {
				ctx.Cancel()
				return
			}
		},
	}})

	if bwGame.Stage() == game.Running {
		pl.Hurt(20, entity.VoidDamageSource{})
	} else {
		bwGame.ForEachActivePlayer(func(pl *player.Player) {
			pl.Message(text.Colourf(language.Translate(pl).Game.JoinGame, database.LobbyNameDisplay.Name(u.Data), len(bwGame.OriginalPlayers()), teamSize*teamCount))
		}, tx)
	}
}

func (h PlayerHandler) HandleQuit(pl *player.Player) {
	delete(core.Players, pl.UUID())

	u := user.GetUser(pl)
	u.Game = nil
	user.Save(pl)
	h.game.RemovePlayerFromTeam(pl)
}

func (h PlayerHandler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.GetUser(pl)

	if listener.CheckChatCoolDown(pl) {
		return
	}

	msgColor := lo.If(u.Data.Rank() <= database.Booster, "white").Else("grey")

	*msg = strings.ReplaceAll(*msg, "§r", "")
	if h.game.Stage() == game.Running {
		oldMsg := *msg
		*msg = text.Colourf("%v<grey>:</grey> <%v>%v</%v>", database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Colour()).Name(u.Data), msgColor, *msg, msgColor)
		if h.game.typeGame == game.TypeBedFight {
			for e := range pl.Tx().Players() {
				p, _ := e.(*player.Player)
				p.Message(*msg)
			}
		} else {
			if strings.HasPrefix(oldMsg, "!") {
				*msg = strings.Replace(*msg, "!", "", 1)
				*msg = text.Colourf("<gold>[SHOUT]</gold> %v<grey>:</grey> <%v>%v</%v>", database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Colour()).Name(u.Data), msgColor, *msg, msgColor)
				for e := range pl.Tx().Players() {
					p, _ := e.(*player.Player)
					p.Message(*msg)
				}
			} else {
				ctx.Cancel()
				h.game.PlayerTeam(pl).ForEachPlayer(pl.Tx(), func(pl *player.Player) {
					pl.Message(text.Colourf(*msg))
				})
				return
			}
		}
	} else {
		*msg = text.Colourf("%v<grey>:</grey> <%v>%v<%v>", database.LobbyNameDisplay.Name(u.Data), msgColor, *msg, msgColor)
		for e := range pl.Tx().Players() {
			p, _ := e.(*player.Player)
			p.Message(*msg)
		}
	}
}

func (PlayerHandler) HandleItemConsume(ctx *player.Context, s item.Stack) {
	if s, ok := s.Item().(item.Potion); ok {
		ctx.Cancel()

		pl := ctx.Val()
		switch s.Type {
		case potion.StrongLeaping():
			pl.AddEffect(effect.New(effect.JumpBoost, 5, 45*time.Second))
		case potion.StrongSwiftness():
			pl.AddEffect(effect.New(effect.Speed, 2, 45*time.Second))
		case potion.LongInvisibility():
			pl.AddEffect(effect.New(effect.Invisibility, 1, 30*time.Second))
		}

		main, off := pl.HeldItems()
		pl.SetHeldItems(main.Grow(-1), off)
	}
}

func (PlayerHandler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	listener.HandleAttackEntity(ctx, e, force, height, critical)

	pl := ctx.Val()
	u := user.GetUser(pl)

	if u.IsCooldownActive(user.Switching, 0, false, false, false) {
		ctx.Cancel()
	}
}

func (h PlayerHandler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	pl := ctx.Val()
	if pl.GameMode() != world.GameModeSpectator && newPos.Y() <= float64(h.game.MapConfig().Void) {
		if h.game.Stage() < game.Running {
			pl.Teleport(h.game.MapConfig().SpawnPoint)
		} else {
			damage := 30.0
			dur := time.Duration(0)
			h.HandleHurt(ctx, &damage, false, &dur, entity.VoidDamageSource{})
		}
	}

	if h.game.Stage() == game.Running {
		distance := float64(h.game.MapConfig().HeightLimit) - pl.Position().Y()
		pl.SendTip(text.Colourf("<dark-red>HEIGHT LIMIT: </dark-red> %v", lo.If(distance <= 0, text.Colourf("<red>REACHED</red>")).Else(text.Colourf("<green>%.1f</green>", distance))))
	}

	if h.game.Stage() == game.Running && h.game.Type() == game.TypeBedWars {
		team := h.game.PlayerTeam(pl)

		if team.Upgrades.HealPool > 0 && utils.Distance(newPos, h.game.MapConfig().TeamSpawnPoints[team.ID()]) <= 20 {
			pl.AddEffect(effect.NewInfinite(effect.Regeneration, 1))
		} else {
			pl.RemoveEffect(effect.Regeneration)
		}

		nearestEnemyTeam := h.game.NearestEnemyTeam(team, pl.Position())
		enemyPos := h.game.MapConfig().BedPositions[nearestEnemyTeam.ID()*2]

		if !h.game.trapIgnore[pl.UUID()] && nearestEnemyTeam.TrapsCount() > 0 && utils.Distance(newPos, enemyPos) <= 10 {
			if nearestEnemyTeam.Upgrades.ActiveTrap == game.None || time.Now().Sub(nearestEnemyTeam.Upgrades.ActivatedSince) > 10*time.Second {
				nearestEnemyTeam.Upgrades.ActiveTrap = nearestEnemyTeam.RemoveTrap()
				nearestEnemyTeam.Upgrades.ActivatedSince = time.Now()
				nearestEnemyTeam.ForEachPlayer(pl.Tx(), func(p *player.Player) {
					p.SendTitle(title.New(text.Colourf(language.Translate(p).BedWars.TrapTriggered)))
				})
			}

			if nearestEnemyTeam.Upgrades.ActiveTrap != game.None {
				switch nearestEnemyTeam.Upgrades.ActiveTrap {
				case game.Regular:
					if _, ok := pl.Effect(effect.Blindness); !ok {
						pl.AddEffect(effect.New(effect.Blindness, 2, 10*time.Second))
						pl.AddEffect(effect.New(effect.Slowness, 1, 10*time.Second))
					}
				case game.CounterOffensive:
					nearestEnemyTeam.ForEachPlayer(pl.Tx(), func(p *player.Player) {
						if _, ok := pl.Effect(effect.Speed); !ok {
							p.AddEffect(effect.New(effect.Speed, 2, 10*time.Second))
							p.AddEffect(effect.New(effect.JumpBoost, 2, 10*time.Second))
						}
					})
				case game.Alarm:
					pl.RemoveEffect(effect.Invisibility)
				case game.MinerFatigue:
					if _, ok := pl.Effect(effect.MiningFatigue); !ok {
						pl.AddEffect(effect.New(effect.MiningFatigue, 1, 10*time.Second))
					}
				default:
					panic("unhandled default case")
				}
			}
		}
	}
}

func (h PlayerHandler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if !listener.HandleHurt(ctx, damage, immune, attackImmunity, src) {
		return
	}

	pl := ctx.Val()
	u := user.GetUser(pl)

	if h.game.Stage() != game.Running {
		ctx.Cancel()
		return
	}

	if _, ok := src.(entity.ExplosionDamageSource); ok {
		*damage *= 0.2
	}

	if _, ok := src.(entity.FallDamageSource); ok && h.game.typeGame == game.TypeBedFight {
		ctx.Cancel()
	}

	if s, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := s.Attacker.(*player.Player); ok {
			if !h.game.EnemyWith(pl, attacker) {
				ctx.Cancel()
				return
			}

			ua := user.GetUser(attacker)
			u.LastHit = attacker.H()
			ua.LastHit = pl.H()
			u.LastHitAt = time.Now()
			ua.LastHitAt = time.Now()

			pl.RemoveEffect(effect.Invisibility)

			if pl.Health() <= *damage {
				onDeath(h.game, pl, u, ua)
				ctx.Cancel()
			}
		}
	} else if u.LastHit != nil && time.Now().Sub(u.LastHitAt) <= 15*time.Second {
		if ea, ok := u.LastHit.Entity(pl.Tx()); ok {
			if pla, ok := ea.(*player.Player); ok && pl.Health() <= *damage {
				onDeath(h.game, pl, u, user.GetUser(pla))
				ctx.Cancel()
				return
			}
		}

		if pl.Health() <= *damage {
			onDeath(h.game, pl, u, nil)
			ctx.Cancel()
		}
	} else if pl.Health() <= *damage {
		onDeath(h.game, pl, u, nil)
		ctx.Cancel()
	}
}

func onDeath(g *BedWars, pl *player.Player, u *user.User, ua *user.User) {
	for _, e := range pl.Effects() {
		if !e.Infinite() {
			pl.RemoveEffect(e.Type())
		}
	}
	pl.Heal(pl.MaxHealth(), effect.InstantHealingSource{})
	pl.SetGameMode(world.GameModeSpectator)
	hadShears := pl.Inventory().ContainsItem(item.NewStack(item.Shears{}, 1))
	newPickaxeTier := pickaxeTier(pl, -1)
	newAxeTier := axeTier(pl, -1)
	inv.CloseContainer(pl)

	finalKill := ""
	if g.PlayerTeam(pl).Status == game.BedBroken {
		finalKill = text.Colourf("<bold><aqua>FINAL KILL!</aqua></bold>")
		g.PlayerTeam(pl).RemovePlayerFromActive(pl)

		if ua != nil {
			ua.GameInfo.BedWars.FinalKills++
			if g.typeGame == game.TypeBedWars {
				ua.Data.Games.BedWars.FinalKills++
			} else {
				ua.Data.Games.BedFight.FinalKills++
			}
		}
	} else {
		go func() {
			i := 3
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				if team := g.PlayerTeam(pl); team != nil {
					if i == 0 {
						pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
							p := e.(*player.Player)
							p.Teleport(g.MapConfig().TeamSpawnPoints[team.ID()])
							if g.typeGame == game.TypeBedFight {
								yaw, _ := living.LookAtExtended(pl.Position(), g.MapConfig().TeamSpawnPoints[1-team.ID()])
								p.Move(mgl64.Vec3{}, yaw-p.Rotation().Yaw(), 0)
								for _, v := range tx.Viewers(p.Position()) {
									v.ViewEntityTeleport(p, p.Position())
								}
							}

							giveKit(p, g)
							if g.typeGame == game.TypeBedWars {
								if hadShears {
									_, _ = u.AddItemWithHBConfig(-1, item.NewStack(item.Shears{}, 1))
								}
								_, _ = u.AddItemWithHBConfig(-1, newPickaxeTier)
								_, _ = u.AddItemWithHBConfig(-1, newAxeTier)
							}

							time.AfterFunc(50*time.Millisecond, func() {
								pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
									e.(*player.Player).SetGameMode(world.GameModeSurvival)
								})
							})
						})
						break
					} else {
						pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.YouDiedTitle)).WithSubtitle(text.Colourf(language.Translate(pl).BedWars.YouDiedSubTitle, i)))
					}
					i--
				} else {
					ticker.Stop()
				}
			}
		}()

		if ua != nil {
			ua.GameInfo.BedWars.Kills++
			if g.typeGame == game.TypeBedWars {
				ua.Data.Games.BedWars.Kills++
			} else {
				ua.Data.Games.BedFight.Kills++
			}
		}
	}

	g.ForEachActivePlayer(func(p *player.Player) {
		if ua == nil {
			p.Message(text.Colourf(language.Translate(p).BedWars.VoidDeath, database.BedWarsNameDisplay(g.PlayerTeam(pl).Colour()).Name(u.Data), finalKill))
		} else {
			c1 := g.PlayerTeam(pl).Colour()
			c2 := g.PlayerTeam(ua.Player()).Colour()
			p.Message(text.Colourf(
				language.Translate(p).BedWars.KilledBy,
				text.Colourf("<%v>%v</%v>", c1, u.Data.Username, c1),
				text.Colourf("<%v>%v</%v>", c2, ua.Data.Username, c2),
				finalKill,
			))
		}
	}, pl.Tx())

	if ua != nil {
		if attacker, ok := ua.Player().H().Entity(pl.Tx()); ok {
			ua.Player().PlaySound(sound.Experience{})
			rewardResources(attacker.(*player.Player), pl)
		}
	}

	pl.Inventory().Clear()

	if g.typeGame == game.TypeBedWars {
		u.Data.Games.BedWars.Deaths++
	} else {
		u.Data.Games.BedFight.Deaths++
	}
}

func (PlayerHandler) HandleHeldSlotChange(ctx *player.Context, from, to int) {
	pl := ctx.Val()
	u := user.GetUser(pl)
	u.IsCooldownActive(user.Switching, time.Duration(core.Config.Pvp.HitRegistration)*time.Millisecond, false, true, false)
}

func (h PlayerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	pl := ctx.Val()
	main, _ := pl.HeldItems()
	if _, ok := main.Item().(item.Bucket); ok {
		time.AfterFunc(50*time.Millisecond, func() {
			pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				pl = e.(*player.Player)
				main, off := pl.HeldItems()
				pl.SetHeldItems(main.Grow(-1), off)
			})
		})
	}
}

func (h PlayerHandler) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	pl := ctx.Val()
	_, ok1 := e.(*ItemsShopVillager)
	_, ok2 := e.(*UpgradesShopVillager)
	if ok1 || ok2 {
		e.(entity.Living).Hurt(0, entity.AttackDamageSource{Attacker: pl})
	}
}

func (h PlayerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	pl := ctx.Val()
	main, off := pl.HeldItems()

	if h.game.Stage() < game.Running {
		ctx.Cancel()
		return
	}

	blocksPlaced[vec3ToString(pos.Vec3())] = b

	if t, ok := b.(block.TNT); ok {
		ctx.Cancel()
		t.Ignite(pos, pl.Tx(), nil)
		pl.SetHeldItems(main.Grow(-1), off)
		return
	}

	for _, v := range slices.Concat(
		h.game.MapConfig().IronGenerators,
		h.game.MapConfig().DiamondGenerators,
		h.game.MapConfig().EmeraldGenerators,
		h.game.MapConfig().ShopVillagerPositions,
		h.game.MapConfig().UpgradesVillagerPositions,
	) {
		if utils.Distance(v, pos.Vec3()) <= 4 {
			ctx.Cancel()
			break
		}
	}

	if h.game.MapConfig().HeightLimit <= pos.Y() {
		ctx.Cancel()
	}
}

func (h PlayerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	pl := ctx.Val()
	u := user.GetUser(pl)

	b := pl.Tx().Block(pos)
	_, isEndstone := b.(block.EndStone)
	_, isPlank := b.(block.Planks)
	bb, isBed := b.(bed.Bed)

	if !isBed && (h.game.typeGame == game.TypeBedWars || !isEndstone && !isPlank) && (h.game.Stage() < game.Running || blocksPlaced[vec3ToString(pos.Vec3())] == nil) {
		pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakMap))
		ctx.Cancel()
	} else {
		blocksPlaced[vec3ToString(pos.Vec3())] = nil
	}

	if isEndstone || isPlank {
		return
	}

	if isBed {
		var teamIndex int
		var bedColor string
		switch bb.Colour {
		case item.ColourRed():
			teamIndex = 0
			bedColor = text.Colourf("<red>Red Bed</red>")
		case item.ColourGreen(), item.ColourLime():
			teamIndex = 1
			bedColor = text.Colourf("<green>Green Bed</green>")
		case item.ColourBlue():
			teamIndex = 2
			bedColor = text.Colourf("<blue>Blue Bed</blue>")
		case item.ColourYellow():
			teamIndex = 3
			bedColor = text.Colourf("<yellow>Yellow Bed</yello>")
		}

		if h.game.PlayerTeam(pl).ID() == teamIndex {
			ctx.Cancel()
			pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakBed))
			return
		}

		h.game.Teams()[teamIndex].Status = game.BedBroken
		h.game.Teams()[teamIndex].ForEachPlayer(pl.Tx(), func(p *player.Player) {
			p.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.BedBreakTitle)).WithSubtitle(text.Colourf(language.Translate(pl).BedWars.BedBreakSubTitle)))
		})

		u.GameInfo.BedWars.BedsBroken++
		if h.game.typeGame == game.TypeBedWars {
			u.Data.Games.BedWars.BedsBroken++
		} else {
			u.Data.Games.BedFight.BedsBroken++
		}

		h.game.playBedBrokenSound(pl.Tx())
		pl.Message(text.Colourf(language.Translate(pl).BedWars.BedBreak, bedColor, database.BedWarsNameDisplay(h.game.PlayerTeam(pl).Colour()).Name(u.Data)))
		return
	}
}

func (PlayerHandler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	ctx.Cancel()
}

func (PlayerHandler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	listener.HandleStartBreak(ctx, pos)
	pl := ctx.Val()
	main, off := pl.HeldItems()

	_, isSword := main.Item().(item.Sword)
	_, isTool := main.Item().(item.Tool)
	if !(isSword || isTool) {
		putItem := func(inv *inventory.Inventory) {
			_, _ = inv.AddItem(main)
			pl.SetHeldItems(item.Stack{}, off)
		}

		if chest, ok := pl.Tx().Block(pos).(block.Chest); ok {
			putItem(chest.Inventory(pl.Tx(), pos))
		} else if _, ok := pl.Tx().Block(pos).(block.EnderChest); ok {
			putItem(pl.EnderChestInventory())
		}
	}
}

func (PlayerHandler) HandlePunchAir(ctx *player.Context) {
	listener.HandlePunchAir(ctx)
}

func (h PlayerHandler) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	pl := ctx.Val()

	if h.game.typeGame == game.TypeBedWars {
		gen := h.game.NearestGenerator(pl.Position(), Iron)
		if gen != nil {
			genPlayers := gen.PlayersWithin(pl.Tx())
			if len(genPlayers) > 1 && !i.Empty() {
				ctx.Cancel()
				split(pl, genPlayers, h)
			}
		}
	}
}

func split(pl *player.Player, genPlayers []*player.Player, h PlayerHandler) {
	genIron := h.game.NearestGenerator(pl.Position(), Iron)
	genGold := h.game.NearestGenerator(pl.Position(), Gold)

	f := func(gen *GeneratorBlockType) {
		gen.UpdateQueue(pl.Tx())
		for _, ent := range gen.ResourcesWithin(pl.Tx()) {
			if be, ok := ent.Behaviour().(*entity.ItemBehaviour); ok {
				if be.Item().Count() > 1 {
					n := be.Item().Count() / len(genPlayers)
					for _, p := range genPlayers {
						pickUp(p, ent, item.NewStack(be.Item().Item(), n), false, pl.Tx())
					}
					utils.Panic(ent.Close())
				} else {
					pickUp(gen.Next(), ent, be.Item(), true, pl.Tx())
				}
				break
			}
		}
	}

	f(genIron)
	f(genGold)
}

func pickUp(pl *player.Player, ent *entity.Ent, stack item.Stack, closeEnt bool, tx *world.Tx) {
	_, _ = pl.Inventory().AddItem(stack)

	collector, ok := world.Entity(pl).(entity.Collector)
	if !ok {
		return
	}

	for _, viewer := range tx.Viewers(ent.Position()) {
		viewer.ViewEntityAction(ent, entity.PickedUpAction{Collector: collector})
	}

	if closeEnt {
		utils.Panic(ent.Close())
	}
}

func (PlayerHandler) HandleItemUse(ctx *player.Context) {
	listener.HandleItemUse(ctx)
}

func vec3ToString(v mgl64.Vec3) string {
	return fmt.Sprintf("(%d, %d, %d)", int(v.X()), int(v.Y()), int(v.Z()))
}

func giveKit(pl *player.Player, g *BedWars) {
	u := user.GetUser(pl)
	t := g.PlayerTeam(pl)

	sword := item.NewStack(item.Sword{Tier: item.ToolTierWood}, 1).AsUnbreakable()
	if t.Upgrades.Sharpness > 0 {
		sword = sword.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, t.Upgrades.Sharpness))
	}
	_, _ = u.AddItemWithHBConfig(0, sword)
	if g.Type() == game.TypeBedFight {
		_, _ = u.AddItemWithHBConfig(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)).AsUnbreakable())
		_, _ = u.AddItemWithHBConfig(2, item.NewStack(item.Axe{Tier: item.ToolTierWood}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)).AsUnbreakable())
		_, _ = u.AddItemWithHBConfig(3, item.NewStack(item.Shears{}, 1).AsUnbreakable())
		_, _ = u.AddItemWithHBConfig(4, item.NewStack(block.Wool{Colour: t.WoolColour()}, 64))
	}

	if len(pl.Armour().Items()) == 0 {
		pl.Armour().Set(
			item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: t.WoolColour().RGBA()}}, 1).AsUnbreakable(),
			item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{Colour: t.WoolColour().RGBA()}}, 1).AsUnbreakable(),
			item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{Colour: t.WoolColour().RGBA()}}, 1).AsUnbreakable(),
			item.NewStack(item.Boots{Tier: item.ArmourTierLeather{Colour: t.WoolColour().RGBA()}}, 1).AsUnbreakable(),
		)

		if t.Upgrades.Protection != 0 {
			for slot, stack := range pl.Armour().Items() {
				utils.Panic(pl.Armour().Inventory().SetItem(slot, stack.WithEnchantments(item.NewEnchantment(enchantment.Protection, t.Upgrades.Protection))))
			}
		}
	}
}

func rewardResources(pl *player.Player, killed *player.Player) {
	var iron, gold, diamond, emerald int
	for _, stack := range killed.Inventory().Items() {
		switch stack.Item().(type) {
		case item.IronIngot:
			iron += stack.Count()
		case item.GoldIngot:
			gold += stack.Count()
		case item.Diamond:
			diamond += stack.Count()
		case item.Emerald:
			emerald += stack.Count()
		}
	}

	if iron > 0 {
		_, _ = pl.Inventory().AddItem(item.NewStack(item.IronIngot{}, iron))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveIron, iron))
	}

	if gold > 0 {
		_, _ = pl.Inventory().AddItem(item.NewStack(item.GoldIngot{}, gold))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveGold, gold))
	}

	if diamond > 0 {
		_, _ = pl.Inventory().AddItem(item.NewStack(item.Diamond{}, diamond))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveDiamond, diamond))
	}

	if emerald > 0 {
		_, _ = pl.Inventory().AddItem(item.NewStack(item.Emerald{}, emerald))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveEmerald, emerald))
	}
}
