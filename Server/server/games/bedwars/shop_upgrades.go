package bedwars

import (
	"server/server/blocks"
	"server/server/game"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type upgradesShop struct {
	game   *BedWars
	team   *game.Team
	player *player.Player
}

func (s *upgradesShop) Main() []item.Stack {
	upgrades := s.team.Upgrades

	items := make([]item.Stack, 54)
	items[10] = s.uShopify(
		item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1),
		"Sharpness",
		"Increases the sword's damage",
		upgrades.Sharpness,
		[]int{8},
		false,
	)
	items[11] = s.uShopify(
		item.NewStack(item.Chestplate{Tier: item.ArmourTierIron{}}, 1),
		"Armour Protection",
		"Increases the armor's protection",
		upgrades.Protection,
		[]int{5, 10, 15, 20},
		false,
	)
	items[12] = s.uShopify(
		item.NewStack(item.Pickaxe{Tier: item.ToolTierGold}, 1),
		"Haste",
		"Increase your tool's haste",
		upgrades.Haste,
		[]int{4, 6},
		false,
	)

	items[14] = s.uShopify(
		item.NewStack(blocks.TripWire{}, 1),
		"Regular Trap",
		"Gives Blindness II & Slowness I effect to the attacker for 8 seconds",
		s.team.TrapsCount(),
		[]int{1, 2, 3},
		true,
	)
	items[15] = s.uShopify(
		item.NewStack(item.Feather{}, 1),
		"Counter-Offensive Trap",
		"Gives Speed II & Jump II to all team members for 15 seconds",
		s.team.TrapsCount(),
		[]int{1, 2, 3},
		true,
	)
	items[16] = s.uShopify(
		item.NewStack(block.Torch{Type: block.SoulFire()}, 1),
		"Alarm Trap",
		"Removes invisibility from infiltrating players",
		s.team.TrapsCount(),
		[]int{1, 2, 3},
		true,
	)

	items[19] = s.uShopify(
		item.NewStack(block.Furnace{}, 1),
		"Forge Upgrade",
		"Increases the speed and cap of resources generated",
		upgrades.GeneratorTier,
		[]int{4, 8, 12, 16},
		false,
	)
	items[20] = s.uShopify(
		item.NewStack(block.Beacon{}, 1),
		"Heal Pool",
		"Gives Regeneration I to all team members whenever inside the team island",
		upgrades.HealPool,
		[]int{3},
		false,
	)
	items[23] = s.uShopify(
		item.NewStack(item.Pickaxe{Tier: item.ToolTierGold}, 1),
		"Miner Fatigue Trap",
		"Decreases the mining speed for enemies entering your base",
		s.team.TrapsCount(),
		[]int{1, 2, 3},
		true,
	)

	for i := 27; i <= 35; i++ {
		items[i] = glassPane.WithCustomName(text.Colourf("<grey>↑ Purchasable</grey>\n<grey>↓ Traps Queue</grey>"))
	}

	noTrap := item.NewStack(block.StainedGlass{Colour: item.ColourGrey()}, 1).WithCustomName(text.Colourf("<red>No Trap!</red>\n<grey>Purchasing a trap will queue it here. Its cost scale based on the number of traps queued</grey>"))

	getTrapItem := func(t game.Trap) item.Stack {
		if t != game.None {
			return items[t.Slot()]
		}
		return noTrap
	}
	items[39] = getTrapItem(upgrades.Traps[0])
	items[40] = getTrapItem(upgrades.Traps[1])
	items[41] = getTrapItem(upgrades.Traps[2])
	return items
}

func (s *upgradesShop) uShopify(stack item.Stack, name, description string, tier int, costs []int, isTrap bool) item.Stack {
	if len(costs)-1 >= tier {
		stack = stack.WithValue("resource", int(Diamond)).WithValue("cost", costs[tier])
	}

	infoMsg := "<yellow>Click to purchase!</yellow>"
	if len(costs)-1 < tier {
		infoMsg = "<emerald>Maxed!</emerald>"
	} else if !canAfford(s.player, stack) {
		infoMsg = "<red>You cannot afford this!</red>"
	}

	lore := []string{
		text.Colourf("<grey>%v</grey>", description),
		"",
	}

	if isTrap {
		if len(costs)-1 >= tier {
			lore = append(lore, text.Colourf("<blue>Cost:</blue> <diamond>%v Diamond%v<diamond>", costs[tier], lo.If(costs[tier] > 1, "s").Else("")))
		}
	} else {
		var costsStr []string
		for t, c := range costs {
			costsStr = append(costsStr, text.Colourf(
				"<%v>Tier %v: %v %v</%v>, <diamond>%v Diamond%v</diamond>",
				lo.If(t < tier, "green").Else("grey"),
				t+1,
				name,
				utils.IntToRoman(t+1),
				lo.If(t < tier, "green").Else("grey"),
				c,
				lo.If(c > 1, "s").Else(""),
			))
		}

		lore = append(lore, costsStr...)
		lore = append(lore, "")
	}
	lore = append(lore, text.Colourf(infoMsg))

	return stack.WithCustomName(text.Colourf("<yellow>%v</yellow>", name)).WithLore(lore...)
}
