package ui

import (
	"server/server/database"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"

	"github.com/bedrock-gophers/inv/inv"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func sendHotBarConfigUI(pl *player.Player) {
	u := user.GetUser(pl)

	menuInv := inventory.New(54, func(slot int, before, after item.Stack) {
		utils.Session(pl).ViewSlotChange(slot, after)
	})
	menu := inv.NewCustomMenu(text.Colourf("<emerald>HotBar Editor</emerald>"), inv.ContainerChest{DoubleChest: true}, menuInv, nil)
	menu.WithStacks(main(pl)...)

	menuInv.Handle(listener.ChestUIHandler{Inventory: menuInv, Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, _ *inventory.Inventory) {
			ctx.Cancel()

			sword, err := menuInv.Item(11)
			if err == nil && sword.Equal(database.Melee.AsStack()) {
				if slot <= 35 {
					go func() {
						pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})

						for i := 11; i <= 35; i++ {
							utils.Panic(menuInv.SetItem(i, item.Stack{}))
						}

						utils.Panic(menuInv.SetItem(13, stack))
					}()
				} else if slot >= 36 && slot <= 44 {
					u.Data.Settings.HotBarConfig[slot-36] = database.None

					go func() {
						pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
						menu.WithStacks(main(pl)...)
					}()
				}
			} else if slot >= 36 && slot <= 44 {
				categoryStack, _ := menuInv.Item(13)
				category := database.HotBarCategoryFromStack(categoryStack)

				for s, c := range u.Data.Settings.HotBarConfig {
					if c == category {
						u.Data.Settings.HotBarConfig[s] = database.None
						break
					}
				}

				u.Data.Settings.HotBarConfig[slot-36] = category

				go func() {
					pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {})
					menu.WithStacks(main(pl)...)
				}()
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

func main(pl *player.Player) []item.Stack {
	u := user.GetUser(pl)

	items := make([]item.Stack, 54)

	items[11] = database.Melee.AsStack()
	items[12] = database.Blocks.AsStack()
	items[13] = database.Bows.AsStack()
	items[14] = database.Potions.AsStack()
	items[15] = database.Utility.AsStack()
	items[21] = database.Shears.AsStack()
	items[22] = database.Pickaxe.AsStack()
	items[23] = database.Axe.AsStack()
	items[31] = database.Ladder.AsStack()

	for i := 0; i < 9; i++ {
		items[i+36] = u.Data.Settings.HotBarConfig[i].AsStack()
	}

	return items
}
