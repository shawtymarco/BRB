package bedwars

import (
	"server/server/game"
	"server/server/inv"
	inv2 "server/server/inv"
	"server/server/living"
	"server/server/user"
	"server/server/utils"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"

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

type ItemsShopVillager struct {
	living.NopLivingType
	*living.Living

	Game *BedWars
	Team *game.Team
}

func NewItemsVillager(pos mgl64.Vec3, name string, game *BedWars, team *game.Team, tx *world.Tx) *ItemsShopVillager {
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
		if pl, ok := src.Attacker.(*player.Player); ok && !user.GetUser(pl).IsCooldownActive(user.Interact, 50*time.Millisecond, false, true, false) {
			SendItemShopUI(&ItemShop{game: v.Game, Player: pl}, false)
		}
	}

	return 0, false
}

func SendItemShopUI(shop *ItemShop, fromLobby bool) {
	pl := shop.Player
	u := user.GetUser(pl)

	menuInv := inventory.New(54, func(slot int, before, after item.Stack) {
		utils.Session(pl).ViewSlotChange(slot, after)
	})
	shop.inv = menuInv
	activeItemShops[pl.UUID()] = shop
	menu := inv.NewCustomMenu(text.Colourf("<emerald>Item Shop</emerald>"), inv.ContainerChest{DoubleChest: true}, menuInv, nil)
	menu.WithStacks(shop.itemShopDashboard(true)...)
	inQuickBuy := true

	var qbSlot int

	menuInv.Handle(inv2.ChestUIHandler{Inventory: menuInv, Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, _ *inventory.Inventory) {
			ctx.Cancel()

			if slot >= 1 && slot <= 7 {
				inQuickBuy = false
				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					menuInv.Clear()
					switch slot {
					case 1:
						menu.WithStacks(shop.Blocks()...)
						return
					case 2:
						menu.WithStacks(shop.Melee()...)
						return
					case 3:
						menu.WithStacks(shop.Armour()...)
						return
					case 4:
						menu.WithStacks(shop.Tools()...)
						return
					case 5:
						menu.WithStacks(shop.Bows()...)
						return
					case 6:
						menu.WithStacks(shop.Potions()...)
						return
					case 7:
						menu.WithStacks(shop.Utility()...)
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
					inQuickBuy = true
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
						shop.IsQuickBuy = true

						qbSlot = slot
						go func() {
							pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
							menuInv.Clear()
							menu.WithStacks(shop.Blocks()...)
						}()
					} else if !fromLobby {
						go func() {
							pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
								owned := strings.Contains(stack.Lore()[len(stack.Lore())-1], "Purchased!")
								if !owned && shop.game.buyItem(pl, stack.AsUnbreakable()) {
									if _, ok := stack.Item().(item.Boots); ok {
										menu.WithStacks(lo.If(inQuickBuy, shop.itemShopDashboard(true)).Else(shop.Armour())...)
									} else if tool, ok := stack.Item().(item.Tool); ok {
										switch tool.ToolType() {
										case item.TypePickaxe:
											utils.Panic(menuInv.SetItem(slot, pickaxeTier(pl, 1)))
										case item.TypeAxe:
											utils.Panic(menuInv.SetItem(slot, axeTier(pl, 1)))
										case item.TypeShears:
											utils.Panic(menuInv.SetItem(slot, axeTier(pl, 1)))
											menu.WithStacks(lo.If(inQuickBuy, shop.itemShopDashboard(true)).Else(shop.Tools())...)
										}
									}
									pl.PlaySound(sound.Experience{})
								} else {
									resource, cost := getCost(stack)
									_ = menuInv.SetItem(slot, shopify(pl, stack, resource, cost, owned, true))
									pl.PlaySound(sound.Deny{})
								}
							})
						}()
					}
				} else if qbMode.Equal(blazePowder) {
					itemId := -1
					for i, s2 := range shop.All() {
						if s2.Equal(stack) {
							itemId = i
							break
						}
					}

					if itemId != -1 {
						shop.IsQuickBuy = false
						u.Data.Settings.QuickBuyConfig[qbSlot] = &itemId
						go func() {
							pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
							menuInv.Clear()
							menu.WithStacks(shop.itemShopDashboard(true)...)
							inQuickBuy = true
						}()
					}
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
