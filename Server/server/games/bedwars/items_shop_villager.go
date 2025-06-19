package bedwars

import (
	"server/server/game"
	"server/server/living"
	"server/server/user"
	"server/server/utils"
	"slices"

	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/bedrock-gophers/inv/inv"
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

type ItemsShopVillager struct {
	living.NopLivingType
	*living.Living

	Game *BedWars
	Team *game.Team
}

func NewShopVillager(pos mgl64.Vec3, name string, game *BedWars, team *game.Team, tx *world.Tx) *ItemsShopVillager {
	m := &ItemsShopVillager{Game: game, Team: team}
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
	l := tx.AddEntity(world.EntitySpawnOpts{NameTag: name, Position: pos}.New(conf.EntityType, conf)).(*ItemsShopVillager)
	return l
}

func (*ItemsShopVillager) EncodeEntity() string {
	return "minecraft:villager"
}
func (*ItemsShopVillager) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
func (v *ItemsShopVillager) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	v.Living = v.NopLivingType.Open(tx, handle, data).(*living.Living)
	return v
}

func (v *ItemsShopVillager) Hurt(dmg float64, src world.DamageSource) (float64, bool) {
	if src, ok := src.(entity.AttackDamageSource); ok {
		if pl, ok := src.Attacker.(*player.Player); ok {
			sendItemShopUI(&itemShop{game: v.Game, team: v.Team, player: pl})
		}
	}

	return 0, false
}

func sendItemShopUI(shop *itemShop) {
	pl := shop.player
	u := user.LookupPlayer(pl)

	menuInv := inventory.New(54, func(slot int, before, after item.Stack) {
		utils.Session(pl).ViewSlotChange(slot, after)
	})
	menu := inv.NewCustomMenu(text.Colourf("<emerald>Item Shop</emerald>"), inv.ContainerChest{DoubleChest: true}, menuInv, nil)
	menu.WithStacks(shop.itemShopDashboard(true)...)

	var qbSlot int

	menuInv.Handle(utils.ChestUIHandler{Inventory: menuInv, Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, _ *inventory.Inventory) {
			ctx.Cancel()

			if slot >= 1 && slot <= 7 {
				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					menuInv.Clear()
					switch slot {
					case 1:
						menu.WithStacks(shop.itemShopBlocks()...)
						return
					case 2:
						menu.WithStacks(shop.itemShopMelee()...)
						return
					case 3:
						menu.WithStacks(shop.itemShopArmour()...)
						return
					case 4:
						menu.WithStacks(shop.itemShopTools()...)
						return
					case 5:
						menu.WithStacks(shop.itemShopBows()...)
						return
					case 6:
						menu.WithStacks(shop.itemShopPotions()...)
						return
					case 7:
						menu.WithStacks(shop.itemShopUtility()...)
						return
					}
				}()
				return
			}

			if stack.Equal(netherStar) {
				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					menuInv.Clear()
					menu.WithStacks(shop.itemShopDashboard(true)...)
				}()
				return
			}

			if stack.Equal(blazeRod) {
				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					_ = menuInv.SetItem(slot, redDye)
				}()
			} else if stack.Equal(redDye) {
				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					_ = menuInv.SetItem(slot, blazeRod)
				}()
			} else if stack.Equal(blazePowder) {
				return
			}

			qbMode, _ := menuInv.Item(53)
			if slices.Contains(quickBuySlots, slot) {
				if qbMode.Equal(blazeRod) {
					if stack.Equal(glassPane) {
						shop.isQuickBuy = true

						qbSlot = slot
						go func() {
							pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
							menuInv.Clear()
							menu.WithStacks(shop.itemShopBlocks()...)
						}()
					} else {
						if shop.game.buyItem(pl, stack) {
							pl.PlaySound(sound.Experience{})
						} else {
							cost, resource := getCost(stack)
							go func() {
								pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
								_ = menuInv.SetItem(slot, shopify(pl, stack, resource, cost, false, true))
							}()
							pl.PlaySound(sound.Deny{})
						}
					}
				} else if qbMode.Equal(blazePowder) {
					shop.isQuickBuy = false

					var itemId int
					for i, s2 := range shop.AllItemShops() {
						if s2.Equal(stack) {
							itemId = i
							break
						}
					}

					u.Data.Settings.QuickBuyConfig[qbSlot] = &itemId
					go func() {
						pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
						menuInv.Clear()
						menu.WithStacks(shop.itemShopDashboard(true)...)
					}()
				} else if qbMode.Equal(redDye) {
					u.Data.Settings.QuickBuyConfig[slot] = nil
					go func() {
						pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
						_ = menuInv.SetItem(slot, glassPane)
					}()
				}
				return
			}
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
