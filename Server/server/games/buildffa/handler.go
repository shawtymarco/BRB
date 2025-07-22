package buildffa

import (
	"fmt"
	"image/color"
	core "server/server"
	"server/server/database"
	"server/server/game"
	"server/server/games/lobby"
	"server/server/items/stacks"
	"server/server/language"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"time"

	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/df-mc/dragonfly/server/item/enchantment"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server/player/title"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var blocksPlaced = make(map[string]*blockPlaced)

type Handler struct {
	player.NopHandler
}

func Join(pl *player.Player, tx *world.Tx) {
	pl.Handle(Handler{})

	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()
	giveKit(pl)

	u := user.GetUser(pl)
	u.Game = Game.Game
	u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BUILD FFA</yellow></bold>"))

	pl.SetNameTag(database.LobbyNameDisplay.Name(u.Data))

	tx.RemoveEntity(pl)

	Game.World().Exec(func(tx *world.Tx) {
		tx.AddEntity(pl.H())
		Game.ForEachActivePlayer(func(pl *player.Player) {
			pl.Message(text.Colourf(language.Translate(pl).BuildFFA.JoinMessage, database.LobbyNameDisplay.Name(u.Data)))
		}, tx)
	})
	Game.AddPlayerToTeam(pl, 1, game.TypeBuildFFA)

	pl.Teleport(Game.MapConfig().SpawnPoint)
}

func (Handler) HandleQuit(pl *player.Player) {
	delete(core.Players, pl.UUID())

	u := user.GetUser(pl)
	u.Game = nil
	user.Save(pl)
	Game.RemovePlayerFromTeam(pl)
}

func (Handler) HandleChat(ctx *player.Context, msg *string) {
	lobby.Handler{}.HandleChat(ctx, msg)
}

func (Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	listener.HandleAttackEntity(ctx, e, force, height, critical)
}

func (h Handler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	pl := ctx.Val()
	if pl.GameMode() != world.GameModeSpectator && newPos.Y() <= float64(Game.MapConfig().Void) {
		damage := 30.0
		immunityDur := time.Duration(0)
		h.HandleHurt(ctx, &damage, false, &immunityDur, entity.VoidDamageSource{})
	}

	distance := float64(Game.MapConfig().HeightLimit) - pl.Position().Y()
	pl.SendTip(text.Colourf("<dark-red>HEIGHT LIMIT: </dark-red> %v", lo.If(distance <= 0, text.Colourf("<red>REACHED</red>")).Else(text.Colourf("<green>%.1f</green>", distance))))
}

func (Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	listener.HandleHurt(ctx, damage, immune, attackImmunity, src)

	pl := ctx.Val()
	u := user.GetUser(pl)

	if _, ok := src.(entity.FallDamageSource); ok {
		ctx.Cancel()
		return
	}

	if u.IsCooldownActive(user.Switching, 500*time.Millisecond, false, false, false) {
		ctx.Cancel()
	}

	if s, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := s.Attacker.(*player.Player); ok {
			ua := user.GetUser(attacker)
			u.LastHit = attacker.H()
			ua.LastHit = pl.H()
			u.LastHitAt = time.Now()
			ua.LastHitAt = time.Now()
			if pl.Health() <= *damage {
				onDeath(pl, u, ua)
				ctx.Cancel()
			}
		}
	} else if u.LastHit != nil && time.Now().Sub(u.LastHitAt) <= 15*time.Second {
		if ea, ok := u.LastHit.Entity(pl.Tx()); ok {
			if pla, ok := ea.(*player.Player); ok && pl.Health() <= *damage {
				onDeath(pl, u, user.GetUser(pla))
				ctx.Cancel()
				return
			}
		}

		if pl.Health() <= *damage {
			onDeath(pl, u, nil)
			ctx.Cancel()
		}
	} else if pl.Health() <= *damage {
		onDeath(pl, u, nil)
		ctx.Cancel()
	}
}

func onDeath(pl *player.Player, u *user.User, ua *user.User) {
	pl.Inventory().Clear()
	pl.Armour().Clear()
	giveKit(pl)
	pl.Heal(20, effect.InstantHealingSource{})
	pl.Teleport(Game.MapConfig().SpawnPoint)

	pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BuildFFA.YouDied)))
	Game.ForEachActivePlayer(func(pl *player.Player) {
		if ua == nil {
			pl.Message(text.Colourf(language.Translate(pl).BuildFFA.VoidDeath, database.LobbyNameDisplay.Name(u.Data)))
		} else {
			pl.Message(text.Colourf(language.Translate(pl).BuildFFA.KilledBy, database.LobbyNameDisplay.Name(u.Data), database.LobbyNameDisplay.Name(ua.Data)))
		}
	}, pl.Tx())

	if ua != nil {
		ua.GameInfo.BuildFFA.Kills++
		ua.Player().Heal(20, effect.InstantHealingSource{})
		ua.Player().PlaySound(sound.Experience{})
	}
}

func (Handler) HandleHeldSlotChange(ctx *player.Context, from, to int) {
	pl := ctx.Val()
	u := user.GetUser(pl)
	u.IsCooldownActive(user.Switching, 500*time.Millisecond, true, true, false)
}

func (Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	pl := ctx.Val()
	h := pl.H()

	if Game.MapConfig().HeightLimit <= pos.Y() {
		ctx.Cancel()
		return
	}

	if w, ok := b.(block.Wool); ok {
		blocksPlaced[vec3ToString(pos.Vec3())] = &blockPlaced{block: pl.Tx().Block(pos), placedAt: time.Now()}
		time.AfterFunc(4*time.Second, func() {
			h.ExecWorld(func(tx *world.Tx, e world.Entity) {
				utils.Panics(e.(*player.Player).Inventory().AddItem(item.NewStack(w, 1)))
			})
		})

		time.AfterFunc(10*time.Second, func() {
			if bp := blocksPlaced[vec3ToString(pos.Vec3())]; bp != nil && time.Now().Sub(bp.placedAt) >= 10*time.Second {
				h.ExecWorld(func(tx *world.Tx, e world.Entity) {
					tx.SetBlock(pos, bp.block, nil)
				})
			}
		})
	}
}

func (Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	pl := ctx.Val()
	if blocksPlaced[vec3ToString(pos.Vec3())] == nil {
		*drops = []item.Stack{}
		b := pl.Tx().Block(pos)
		time.AfterFunc(10*time.Second, func() {
			pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				tx.SetBlock(pos, b, nil)
			})
		})
	} else {
		*drops = []item.Stack{}
		blocksPlaced[vec3ToString(pos.Vec3())] = nil
	}
}

func (Handler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	ctx.Cancel()
}

func (Handler) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if _, ok := src.(effect.RegenerationHealingSource); ok {
		ctx.Cancel()
	}
}

func (Handler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	listener.HandleStartBreak(ctx, pos)
}

func (Handler) HandlePunchAir(ctx *player.Context) {
	listener.HandlePunchAir(ctx)
}

func (Handler) HandleItemUse(ctx *player.Context) {
	listener.HandleItemUse(ctx)
}

func giveKit(pl *player.Player) {
	u := user.GetUser(pl)
	utils.Panics(u.AddItemWithHBConfig(0, item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1).AsUnbreakable()))
	utils.Panics(u.AddItemWithHBConfig(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)).AsUnbreakable()))
	utils.Panics(u.AddItemWithHBConfig(2, item.NewStack(item.Shears{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)).AsUnbreakable()))

	utils.Panics(u.AddItemWithHBConfig(4, item.NewStack(block.Wool{Colour: item.ColourLime()}, 64)))
	utils.Panics(u.AddItemWithHBConfig(8, item.NewStack(stacks.KnockBackStick{}, 1).AsUnbreakable().WithEnchantments(item.NewEnchantment(stacks.CustomKnockBack{}, 1)).WithCustomName(text.Colourf("<green>Knockback Stick</green>"))))

	pl.Armour().Set(
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: color.RGBA{G: 255}}}, 1).AsUnbreakable(),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{Colour: color.RGBA{G: 255}}}, 1).AsUnbreakable(),
		item.NewStack(item.Leggings{Tier: item.ArmourTierIron{}}, 1).AsUnbreakable(),
		item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1).AsUnbreakable(),
	)
}

func vec3ToString(v mgl64.Vec3) string {
	return fmt.Sprintf("(%d, %d, %d)", int(v.X()), int(v.Y()), int(v.Z()))
}

type blockPlaced struct {
	block    world.Block
	placedAt time.Time
}
