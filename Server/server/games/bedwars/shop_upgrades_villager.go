package bedwars

import (
	"server/server/game"
	"server/server/inv"
	"server/server/living"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item/enchantment"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type UpgradesShopVillager struct {
	living.NopLivingType
	*living.Living

	Game *BedWars
	Team *game.Team
}

func NewUpgradesVillager(pos mgl64.Vec3, name string, game *BedWars, team *game.Team, tx *world.Tx) *UpgradesShopVillager {
	m := &UpgradesShopVillager{Game: game, Team: team}
	conf := living.Config{
		EntityType: m,
		MaxHealth:  1,
		Speed:      0,
		Drops:      []living.Drop{},
		MovementComputer: &entity.MovementComputer{
			Gravity:           0,
			Drag:              0,
			DragBeforeGravity: false,
		},
		EyeHeight: 0.3,
	}
	l := tx.AddEntity(world.EntitySpawnOpts{NameTag: name, Position: pos}.New(conf.EntityType, conf)).(*UpgradesShopVillager)
	return l
}

func (*UpgradesShopVillager) EncodeEntity() string {
	return "minecraft:villager"
}
func (*UpgradesShopVillager) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
func (v *UpgradesShopVillager) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	v.Living = v.NopLivingType.Open(tx, handle, data).(*living.Living)
	return v
}

func (v *UpgradesShopVillager) Hurt(dmg float64, src world.DamageSource) (float64, bool) {
	if src, ok := src.(entity.AttackDamageSource); ok {
		if pl, ok := src.Attacker.(*player.Player); ok {
			sendUpgradesShopUI(&upgradesShop{game: v.Game, team: v.Team, player: pl})
		}
	}

	return 0, false
}

func sendUpgradesShopUI(shop *upgradesShop) {
	pl := shop.player

	menuInv := inventory.New(54, func(slot int, before, after item.Stack) {
		utils.Session(pl).ViewSlotChange(slot, after)
	})
	menu := inv.NewCustomMenu(text.Colourf("<emerald>Upgrades Shop</emerald>"), inv.ContainerChest{DoubleChest: true}, menuInv, nil)
	menu.WithStacks(shop.Main()...)

	menuInv.Handle(inv.ChestUIHandler{Inventory: menuInv, Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, _ *inventory.Inventory) {
			ctx.Cancel()

			if slot > 26 {
				return
			}

			go func() {
				pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
					if !strings.Contains(stack.Lore()[len(stack.Lore())-1], "Purchased!") && shop.game.buyUpgrade(pl, stack) {
						switch slot {
						case 10:
							shop.team.Upgrades.Sharpness++
							pl.PlaySound(sound.Experience{})

							shop.team.ForEachPlayer(tx, func(p *player.Player) {
								for invSlot, invStack := range p.Inventory().Items() {
									if _, ok := invStack.Item().(item.Sword); ok {
										utils.Panic(p.Inventory().SetItem(invSlot, invStack.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, shop.team.Upgrades.Sharpness))))
									}
								}
							})
						case 11:
							shop.team.Upgrades.Protection++
							pl.PlaySound(sound.Experience{})
							shop.team.ForEachPlayer(tx, func(p *player.Player) {
								for invSlot, invStack := range p.Armour().Items() {
									utils.Panic(p.Armour().Inventory().SetItem(invSlot, invStack.WithEnchantments(item.NewEnchantment(enchantment.Protection, shop.team.Upgrades.Protection))))
								}
							})
						case 12:
							shop.team.Upgrades.Haste++
							pl.PlaySound(sound.Experience{})
							shop.team.ForEachPlayer(tx, func(p *player.Player) {
								p.AddEffect(effect.NewInfinite(effect.Haste, shop.team.Upgrades.Haste))
							})
						case 19:
							shop.team.Upgrades.GeneratorTier++
							pl.PlaySound(sound.Experience{})
							if shop.team.Upgrades.GeneratorTier == 3 {
								(&GeneratorSettings{
									Active:    true,
									Game:      shop.game,
									Resource:  Emerald,
									Tier:      0,
									Name:      "",
									Cap:       4,
									SpawnRate: 30 * time.Second,
								}).New(shop.game.MapConfig().IronGenerators[shop.team.ID()], tx)
							}
						case 20:
							shop.team.Upgrades.HealPool++
							pl.PlaySound(sound.Experience{})
						case 14:
							if !shop.team.IsTrapsFull() {
								shop.team.AddTrap(game.Regular)
								pl.PlaySound(sound.Experience{})
							} else {
								pl.PlaySound(sound.Deny{})
							}
						case 15:
							if !shop.team.IsTrapsFull() {
								shop.team.AddTrap(game.CounterOffensive)
								pl.PlaySound(sound.Experience{})
							} else {
								pl.PlaySound(sound.Deny{})
							}
						case 16:
							if !shop.team.IsTrapsFull() {
								shop.team.AddTrap(game.Alarm)
								pl.PlaySound(sound.Experience{})
							} else {
								pl.PlaySound(sound.Deny{})
							}
						case 23:
							if !shop.team.IsTrapsFull() {
								shop.team.AddTrap(game.MinerFatigue)
								pl.PlaySound(sound.Experience{})
							} else {
								pl.PlaySound(sound.Deny{})
							}
						}
					} else {
						pl.PlaySound(sound.Deny{})
					}

					menu.WithStacks(shop.Main()...)
				})
			}()
		},
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			ctx.Cancel()
			pl.PlaySound(sound.Deny{})
		},
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			ctx.Cancel()
			pl.PlaySound(sound.Deny{})
		},
	}})

	inv.SendMenu(pl, menu)
}
