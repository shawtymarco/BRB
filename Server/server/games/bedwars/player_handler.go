package bedwars

import (
	"fmt"
	"server/server/blocks/bed"
	"server/server/database"
	"server/server/game"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"slices"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/item/inventory"

	"github.com/bedrock-gophers/inv/inv"

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
	"github.com/df-mc/dragonfly/server/player/chat"
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

	pl.Handle(PlayerHandler{game: bwGame})

	tx.RemoveEntity(pl)
	bwGame.World().Exec(func(tx *world.Tx) {
		tx.AddEntity(pl.H())
	})

	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()

	u := user.LookupPlayer(pl)
	u.Game = bwGame.Game
	switch typeGame {
	case game.TypeBedWars:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDWARS</yellow></bold>"))
	case game.TypeBedFight:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDFIGHT</yellow></bold>"))
	default:
		panic("Unhandled game type")
	}

	bwGame.AddPlayerToTeam(pl, teamSize)
	pl.Teleport(bwGame.MapConfig().SpawnPoint)

	if bwGame.Stage() == game.Running {
		pl.Hurt(20, entity.VoidDamageSource{})
	} else {
		bwGame.ForEachActivePlayer(func(pl *player.Player) {
			pl.Message(text.Colourf(language.Translate(pl).Game.JoinGame, database.LobbyNameDisplay.Name(u.Data), len(bwGame.OriginalPlayers()), teamSize*teamCount))
		})
	}
}

func (h PlayerHandler) HandleQuit(pl *player.Player) {
	u := user.LookupPlayer(pl)
	u.Game = nil
	//user.Save(pl) // TODO: Uncomment (Bug from bots)
	h.game.RemovePlayerFromTeam(pl)
	lobby.Join(pl)
}

func (h PlayerHandler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	if listener.CheckChatCoolDown(pl) {
		return
	}

	*msg = strings.ReplaceAll(*msg, "§r", "")
	var newMsg string
	if h.game.Stage() == game.Running {
		if strings.HasPrefix(*msg, "!") {
			*msg = strings.Replace(*msg, "!", "", 1)
			newMsg = fmt.Sprintf("%v<white>: %v</white>", database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Colour()).Name(u.Data), *msg)
		} else {
			ctx.Cancel()
			newMsg = fmt.Sprintf("<grey>[TEAM CHAT]</grey> %v<white>: %v</white>", database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Colour()).Name(u.Data), *msg)
			h.game.PlayerTeam(pl).ForEachPlayer(pl.Tx(), func(pl *player.Player) {
				pl.Message(text.Colourf(newMsg))
			})
			return
		}
	} else {
		newMsg = fmt.Sprintf("%v<white>: %v</white>", database.LobbyNameDisplay.Name(u.Data), *msg)
	}
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func (PlayerHandler) HandleItemConsume(ctx *player.Context, s item.Stack) {
	if s, ok := s.Item().(item.Potion); ok {
		pl := ctx.Val()
		switch s.Type {
		case potion.StrongLeaping():
			pl.AddEffect(effect.New(effect.JumpBoost, 5, 45*time.Second))
		case potion.StrongSwiftness():
			pl.AddEffect(effect.New(effect.Speed, 2, 45*time.Second))
		case potion.LongInvisibility():
			pl.AddEffect(effect.New(effect.Invisibility, 1, 30*time.Second))
		}
	}
}

func (PlayerHandler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	listener.HandleAttackEntity(ctx, e, force, height, critical)
}

func (h PlayerHandler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	pl := ctx.Val()
	if newPos.Y() <= float64(h.game.MapConfig().Void) {
		if h.game.Stage() < game.Running {
			pl.Teleport(h.game.MapConfig().SpawnPoint)
		} else {
			damage := 30.0
			dur := time.Duration(0)
			h.HandleHurt(ctx, &damage, false, &dur, entity.VoidDamageSource{})
		}
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
							p.AddEffect(effect.New(effect.Speed, 2, 15*time.Second))
							p.AddEffect(effect.New(effect.JumpBoost, 2, 15*time.Second))
						}
					})
				case game.Alarm:
					pl.RemoveEffect(effect.Invisibility)
				case game.MinerFatigue:
					if _, ok := pl.Effect(effect.MiningFatigue); !ok {
						pl.AddEffect(effect.New(effect.MiningFatigue, 3, 10*time.Second))
					}
				default:
					panic("unhandled default case")
				}
			}
		}
	}
}

func (h PlayerHandler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	listener.HandleHurt(ctx, damage, immune, attackImmunity, src)

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	if h.game.Stage() < game.Running {
		ctx.Cancel()
		return
	}

	if _, ok := src.(entity.ExplosionDamageSource); ok {
		*damage *= 0.2
	}

	if s, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := s.Attacker.(*player.Player); ok {
			ua := user.LookupPlayer(attacker)
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
				onDeath(h.game, pl, u, user.LookupPlayer(pla))
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
		if g.pickaxeTierPlayers[pl.UUID()] > 1 {
			g.pickaxeTierPlayers[pl.UUID()]--
		}
		if g.axeTierPlayers[pl.UUID()] > 1 {
			g.axeTierPlayers[pl.UUID()]--
		}

		h := pl.H()
		go func() {
			i := 5
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				if team := g.PlayerTeam(pl); team != nil {
					if i == 0 {
						h.ExecWorld(func(tx *world.Tx, e world.Entity) {
							pl := e.(*player.Player)
							pl.Teleport(g.MapConfig().TeamSpawnPoints[team.ID()])
							pl.SetGameMode(world.GameModeSurvival)
							giveKit(pl, g)
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

	ha := ua.Player().H()
	go func() {
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
		})

		if ua != nil {
			pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				if attacker, ok := ha.Entity(tx); ok {
					ua.Player().PlaySound(sound.Experience{})
					rewardResources(attacker.(*player.Player), e.(*player.Player))
				}
				e.(*player.Player).Armour().Clear()
			})
		}

	}()

	if g.typeGame == game.TypeBedWars {
		u.Data.Games.BedWars.Deaths++
	} else {
		u.Data.Games.BedFight.Deaths++
	}
}

func (h PlayerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	pl := ctx.Val()
	main, _ := pl.HeldItems()

	if b, ok := main.Item().(item.Bucket); ok && b.Content == item.LiquidBucketContent(block.Water{}) {
		h := pl.H()
		time.AfterFunc(50*time.Millisecond, func() {
			h.ExecWorld(func(tx *world.Tx, e world.Entity) {
				pl = e.(*player.Player)
				main, off := pl.HeldItems()
				pl.SetHeldItems(main.Grow(-1), off)
			})
		})
	}
}

func (h PlayerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	pl := ctx.Val()
	main, off := pl.HeldItems()
	if h.game.Stage() < game.Running {
		pl.Message(text.Colourf(language.Translate(pl).Game.Error.CannotBreakBlocksBecauseGameNotStarted))
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
		if utils.Distance(v, pos.Vec3()) <= 3 {
			ctx.Cancel()
			break
		}
	}
}

func (h PlayerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	pl := ctx.Val()
	u := user.LookupPlayer(pl)
	b := pl.Tx().Block(pos)
	_, isEndstone := b.(block.EndStone)
	_, isPlank := b.(block.Planks)

	if isEndstone || isPlank {
		return
	}

	if bb, isBed := b.(bed.Bed); isBed {
		var teamIndex int
		var bedColor string
		switch bb.Colour {
		case item.ColourRed():
			teamIndex = 0
			bedColor = text.Colourf("<red>Red</red> Bed")
		case item.ColourBlue():
			teamIndex = 1
			bedColor = text.Colourf("<blue>Blue</blue> Bed")
		case item.ColourGreen():
			teamIndex = 2
			bedColor = text.Colourf("<green>Green</green> Bed")
		case item.ColourYellow():
			teamIndex = 3
			bedColor = text.Colourf("<yellow>Yellow</yello> Bed")
		}

		if h.game.PlayerTeam(pl).ID() == teamIndex {
			ctx.Cancel()
			pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakBed))
			return
		}

		h.game.Teams()[teamIndex].Status = game.BedBroken

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

	if blocksPlaced[vec3ToString(pos.Vec3())] == nil {
		pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakMap))
		ctx.Cancel()
	} else {
		blocksPlaced[vec3ToString(pos.Vec3())] = nil
	}
}

func (PlayerHandler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	ctx.Cancel()
}

func (PlayerHandler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	listener.HandleStartBreak(ctx, pos)
	pl := ctx.Val()
	main, off := pl.HeldItems()

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

func (PlayerHandler) HandlePunchAir(ctx *player.Context) {
	listener.HandlePunchAir(ctx)
}

func (h PlayerHandler) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	pl := ctx.Val()
	gen := h.game.NearestGenerator(pl.Position(), Iron)

	genPlayers := gen.PlayersWithin(pl.Tx())
	if len(genPlayers) > 1 && !i.Empty() {
		ctx.Cancel()
		split(pl, genPlayers, h)
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
	u := user.LookupPlayer(pl)
	t := g.PlayerTeam(pl)

	sword := item.NewStack(item.Sword{Tier: item.ToolTierWood}, 1).AsUnbreakable()
	if t.Upgrades.Sharpness > 0 {
		sword = sword.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, t.Upgrades.Sharpness))
	}
	_, _ = u.AddItemWithHBConfig(0, sword)
	if g.Type() == game.TypeBedFight {
		_, _ = u.AddItemWithHBConfig(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1).AsUnbreakable())
		_, _ = u.AddItemWithHBConfig(2, item.NewStack(item.Axe{Tier: item.ToolTierWood}, 1).AsUnbreakable())
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

	if g.pickaxeTierPlayers[pl.UUID()] != 0 {
		_, _ = u.AddItemWithHBConfig(-1, pickaxeTier(pl, g.pickaxeTierPlayers[pl.UUID()]))
	}
	if g.axeTierPlayers[pl.UUID()] != 0 {
		_, _ = u.AddItemWithHBConfig(-1, axeTier(pl, g.pickaxeTierPlayers[pl.UUID()]))
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
		utils.Panics(pl.Inventory().AddItem(item.NewStack(item.IronIngot{}, iron)))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveIron, iron))
	}

	if gold > 0 {
		utils.Panics(pl.Inventory().AddItem(item.NewStack(item.GoldIngot{}, gold)))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveGold, gold))
	}

	if diamond > 0 {
		utils.Panics(pl.Inventory().AddItem(item.NewStack(item.Diamond{}, diamond)))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveDiamond, diamond))
	}

	if emerald > 0 {
		utils.Panics(pl.Inventory().AddItem(item.NewStack(item.Emerald{}, emerald)))
		pl.Message(text.Colourf(language.Translate(pl).BedWars.GiveEmerald, emerald))
	}
}
