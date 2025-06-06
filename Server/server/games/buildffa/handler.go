package buildffa

import (
	"fmt"
	"image/color"
	"server/server/database"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"strings"
	"time"

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
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var blocksPlaced = make(map[string]time.Time)

type Handler struct {
	player.NopHandler
}

func Join(pl *player.Player, tx *world.Tx) {
	pl.Handle(Handler{})

	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()
	giveKit(pl)

	u := user.LookupPlayer(pl)
	u.Game = Game.Game
	u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BUILDFFA</yellow></bold>"))

	pl.SetNameTag(database.LobbyNameDisplay.Name(u.Data))

	Game.ForEachPlayer(func(pl *player.Player) {
		pl.Message(text.Colourf(language.Translate(pl).BuildFFA.JoinMessage, database.LobbyNameDisplay.Name(u.Data)))
	})

	tx.RemoveEntity(pl)
	Game.World().Exec(func(tx *world.Tx) {
		tx.AddEntity(pl.H())
	})
	Game.AddPlayerToTeam(pl, 1)

	pl.Teleport(Game.MapConfig().SpawnPoint)
}

func (Handler) HandleQuit(pl *player.Player) {
	user.Save(pl)
	Game.RemovePlayerFromTeam(pl)
	lobby.Join(pl)
}

func (Handler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v<white>: %v</white>", database.LobbyNameDisplay.Name(u.Data), *msg)
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func (Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	listener.HandleAttackEntity(ctx, e, force, height, critical)
}

func (Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	listener.HandleHurt(ctx, damage, immune, attackImmunity, src)

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	if s, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := s.Attacker.(*player.Player); ok {
			ua := user.LookupPlayer(attacker)
			if !immune && pl.Health() <= *damage {
				ua.GameInfo.BuildFFA.Kills++
				onDeath(pl, u, ua)

				ctx.Cancel()
			}
		}
	} else if _, ok := src.(entity.VoidDamageSource); ok {
		onDeath(pl, u, nil)
		ctx.Cancel()
	} else {
		ctx.Cancel()
	}
}

func onDeath(pl *player.Player, u *user.User, ua *user.User) {
	pl.Heal(20, effect.InstantHealingSource{})
	pl.Teleport(Game.MapConfig().SpawnPoint)

	pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BuildFFA.YouDied)))
	Game.ForEachPlayer(func(pl *player.Player) {
		if ua == nil {
			pl.Message(text.Colourf(language.Translate(pl).BuildFFA.VoidDeath, database.LobbyNameDisplay.Name(u.Data)))
		} else {
			pl.Message(text.Colourf(language.Translate(pl).BuildFFA.KilledBy, database.LobbyNameDisplay.Name(u.Data), database.LobbyNameDisplay.Name(ua.Data)))
		}
	})
}

func (Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	pl := ctx.Val()
	h := pl.H()
	if w, ok := b.(block.Wool); ok {
		blocksPlaced[vec3ToString(pos.Vec3())] = time.Now()
		time.AfterFunc(4*time.Second, func() {
			h.ExecWorld(func(tx *world.Tx, e world.Entity) {
				utils.Panics(e.(*player.Player).Inventory().AddItem(item.NewStack(w, 1)))
			})
		})

		time.AfterFunc(10*time.Second, func() {
			if !blocksPlaced[vec3ToString(pos.Vec3())].IsZero() && time.Now().Sub(blocksPlaced[vec3ToString(pos.Vec3())]) >= 10*time.Second {
				h.ExecWorld(func(tx *world.Tx, e world.Entity) {
					tx.SetBlock(pos, block.Air{}, nil)
				})
			}
		})
	}
}

func (Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if blocksPlaced[vec3ToString(pos.Vec3())].IsZero() {
		ctx.Cancel()
	} else {
		*drops = []item.Stack{}
		blocksPlaced[vec3ToString(pos.Vec3())] = time.Time{}
	}
}

func (Handler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	ctx.Cancel()
}

func (Handler) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	ctx.Cancel()
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
	utils.Panic(pl.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1).AsUnbreakable()))
	utils.Panic(pl.Inventory().SetItem(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1).AsUnbreakable()))
	utils.Panic(pl.Inventory().SetItem(2, item.NewStack(item.Shears{}, 1).AsUnbreakable()))

	utils.Panic(pl.Inventory().SetItem(4, item.NewStack(block.Wool{Colour: item.ColourGreen()}, 64)))
	utils.Panic(pl.Inventory().SetItem(8, item.NewStack(KnockBackStick{}, 1).AsUnbreakable().WithEnchantments(item.NewEnchantment(enchantment.Knockback, 2)).WithCustomName(text.Colourf("<green>Knockback Stick</green>"))))

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
