package bedwars

import (
	"server/server/game"
	"server/server/items/stacks"
	"server/server/user"
	"strings"

	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/google/uuid"

	"github.com/df-mc/dragonfly/server/item/potion"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/item"
)

var activeItemShops = make(map[uuid.UUID]*ItemShop)

var quickBuySlots = []int{
	19, 20, 21, 22, 23, 24, 25,
	28, 29, 30, 31, 32, 33, 34,
	37, 38, 39, 40, 41, 42, 43,
}

var blazeRod = item.NewStack(item.BlazeRod{}, 1).WithCustomName(text.Colourf("<yellow>Remove an item</yellow>\n<grey>Click here to remove any items from the Quick Buy slot</grey>"))
var blazePowder = item.NewStack(item.BlazePowder{}, 1).WithCustomName(text.Colourf("<yellow>Select an item</yellow>\n<grey>Select any item you want to occupy the Quick Buy slot you just selected</grey>"))
var redDye = item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).WithCustomName(text.Colourf("<red>Removing an item</red>\n<red>Select any item to remove them from a Quick Buy slot</red>"))
var glassPane = item.NewStack(block.StainedGlassPane{Colour: item.ColourGrey()}, 1).WithCustomName(" ")
var netherStar = item.NewStack(item.NetherStar{}, 1).WithCustomName(text.Colourf("<aqua>View QuickBuy</aqua>"))

type ItemShop struct {
	inv    *inventory.Inventory
	game   *BedWars
	team   *game.Team
	Player *player.Player

	IsQuickBuy bool
}

func (s *ItemShop) itemShopDashboard(showQuickBuy bool) []item.Stack {
	items := make([]item.Stack, 54)
	items[1] = item.NewStack(block.Concrete{}, 1).WithCustomName(text.Colourf("<aqua>Blocks</aqua>"))
	items[2] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).WithCustomName(text.Colourf("<aqua>Melee</aqua>"))
	items[3] = item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Armour</aqua>"))
	items[4] = item.NewStack(item.Pickaxe{Tier: item.ToolTierStone}, 1).WithCustomName(text.Colourf("<aqua>Tools</aqua>"))
	items[5] = item.NewStack(item.Bow{}, 1).WithCustomName(text.Colourf("<aqua>Bows</aqua>"))
	items[6] = item.NewStack(block.BrewingStand{}, 1).WithCustomName(text.Colourf("<aqua>Potions</aqua>"))
	items[7] = item.NewStack(block.TNT{}, 1).WithCustomName(text.Colourf("<aqua>Utility</aqua>"))

	if showQuickBuy {
		u := user.GetUser(s.Player)
		allShops := s.All()
		for _, slot := range quickBuySlots {
			if id := u.Data.Settings.QuickBuyConfig[slot]; id != nil {
				items[slot] = allShops[*id]
			} else {
				items[slot] = glassPane
			}
		}
	} else {
		for _, slot := range quickBuySlots {
			items[slot] = item.NewStack(block.Air{}, 1)
		}
	}

	items[49] = netherStar
	items[53] = lo.If(s.IsQuickBuy, blazePowder).Else(blazeRod)
	return items
}

func (s *ItemShop) Blocks() []item.Stack {
	items := s.itemShopDashboard(false)
	if s.team != nil {
		items[19] = shopify(s.Player, item.NewStack(block.Wool{Colour: s.team.WoolColour()}, 16), Iron, 4, false, false)
		items[20] = shopify(s.Player, item.NewStack(block.StainedGlass{Colour: s.team.WoolColour()}, 4), Iron, 8, false, false)
	} else {
		items[19] = shopify(s.Player, item.NewStack(block.Wool{Colour: item.ColourWhite()}, 16), Iron, 4, false, false)
		items[20] = shopify(s.Player, item.NewStack(block.StainedGlass{Colour: item.ColourWhite()}, 4), Iron, 8, false, false)
	}
	items[21] = shopify(s.Player, item.NewStack(block.EndStone{}, 12), Iron, 24, false, false)
	items[22] = shopify(s.Player, item.NewStack(block.Ladder{}, 16), Iron, 4, false, false)
	items[23] = shopify(s.Player, item.NewStack(block.Planks{Wood: lo.If(s.Player != nil, user.GetUser(s.Player).Data.Cosmetics.SelectedWoodType).Else(block.OakWood())}, 16), Gold, 4, false, false)
	items[24] = shopify(s.Player, item.NewStack(block.Obsidian{}, 4), Gold, 3, false, false)
	return items
}

func (s *ItemShop) Melee() []item.Stack {
	items := s.itemShopDashboard(false)
	i19 := shopify(s.Player, item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1), Iron, 10, false, false)
	i20 := shopify(s.Player, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1), Gold, 7, false, false)
	i21 := shopify(s.Player, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1), Emerald, 3, false, false)

	if s.team != nil && s.team.Upgrades.Sharpness > 0 {
		i19 = i19.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
		i20 = i20.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
		i21 = i21.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
	}

	items[19] = i19
	items[20] = i20
	items[21] = i21

	items[22] = shopify(s.Player, item.NewStack(stacks.KnockBackStick{}, 1).AsUnbreakable().WithEnchantments(item.NewEnchantment(stacks.CustomKnockBack{}, 1)).WithCustomName(text.Colourf("<green>Knockback Stick</green>")), Gold, 5, false, false)
	return items
}

func (s *ItemShop) Armour() []item.Stack {
	var chain, iron, diamond bool
	if boots, ok := s.Player.Armour().Boots().Item().(item.Boots); ok {
		t := boots.Tier
		chain = t == item.ArmourTierChain{}
		iron = t == item.ArmourTierIron{}
		diamond = t == item.ArmourTierDiamond{}
	}
	items := s.itemShopDashboard(false)

	i19 := shopify(s.Player, item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1), Iron, 30, chain || iron || diamond, false)
	i20 := shopify(s.Player, item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1), Gold, 12, iron || diamond, false)
	i21 := shopify(s.Player, item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1), Emerald, 6, diamond, false)

	if s.team != nil && s.team.Upgrades.Protection != 0 {
		i19 = i19.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
		i20 = i20.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
		i21 = i21.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
	}

	items[19] = i19
	items[20] = i20
	items[21] = i21
	return items
}

func (s *ItemShop) Tools() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.Player, item.NewStack(item.Shears{}, 1), Iron, 10, s.Player.Inventory().ContainsItemFunc(1, func(stack item.Stack) bool {
		_, isShears := stack.Item().(item.Shears)
		return isShears
	}), false)

	if s.game != nil {
		items[20] = pickaxeTier(s.Player, 1)
		items[21] = axeTier(s.Player, 1)
	} else {
		items[20] = pickaxeTier(s.Player, 1)
		items[21] = axeTier(s.Player, 1)
	}
	return items
}

func (s *ItemShop) Bows() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.Player, item.NewStack(item.Arrow{}, 8), Gold, 2, false, false)
	items[20] = shopify(s.Player, item.NewStack(item.Bow{}, 1), Gold, 12, false, false)
	items[21] = shopify(s.Player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 1)), Gold, 20, false, false)
	items[22] = shopify(s.Player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 2), item.NewEnchantment(enchantment.Punch, 1)), Emerald, 4, false, false)
	return items
}

func (s *ItemShop) Potions() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.Player, item.NewStack(item.Potion{Type: potion.StrongLeaping()}, 1), Emerald, 1, false, false)
	items[20] = shopify(s.Player, item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1), Emerald, 1, false, false)
	items[21] = shopify(s.Player, item.NewStack(item.Potion{Type: potion.LongInvisibility()}, 1), Emerald, 2, false, false)
	return items
}

func (s *ItemShop) Utility() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.Player, item.NewStack(item.GoldenApple{}, 1), Gold, 3, false, false)
	items[20] = shopify(s.Player, item.NewStack(SilverfishSnowball{game: s.game}, 1), Iron, 20, false, false)
	items[21] = editName(shopify(s.Player, item.NewStack(block.TNT{}, 1), Gold, 8, false, false), "TNT")
	items[22] = shopify(s.Player, item.NewStack(item.EnderPearl{}, 1), Emerald, 4, false, false)
	items[23] = shopify(s.Player, item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, 1), Gold, 3, false, false)
	if s.team != nil {
		items[24] = editName(shopify(s.Player, item.NewStack(BridgeEgg{Block: block.Wool{Colour: s.team.WoolColour()}}, 1), Emerald, 1, false, false), text.Colourf("<green>Bridge Egg</green>"))
	} else {
		items[24] = editName(shopify(s.Player, item.NewStack(BridgeEgg{Block: block.Wool{Colour: item.ColourWhite()}}, 1), Emerald, 1, false, false), text.Colourf("<green>Bridge Egg</green>"))
	}
	items[25] = editName(shopify(s.Player, item.NewStack(NewMagicMilk(s.game), 1), Gold, 3, false, false), text.Colourf("<green>Magic Milk</green>"))
	items[28] = shopify(s.Player, item.NewStack(block.Sponge{}, 4), Gold, 3, false, false)
	items[29] = shopify(s.Player, item.NewStack(BedCompass{BedWars: s.game}, 1), Emerald, 2, false, false)
	return items
}

func (s *ItemShop) All() []item.Stack {
	var res []item.Stack
	for slot, i := range s.Blocks() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Melee() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Armour() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Tools() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Bows() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Potions() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.Utility() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}

	return res
}

func shopify(pl *player.Player, s item.Stack, resource Resource, cost int, owned bool, fullInv bool) item.Stack {
	s = s.WithValue("resource", int(resource)).WithValue("cost", cost)

	infoMsg := "<yellow>Click to purchase!</yellow>"
	if owned {
		infoMsg = "<emerald>Purchased!</emerald>"
	} else if !canAfford(pl, s) {
		infoMsg = "<red>You cannot afford this!</red>"
	} else if fullInv {
		infoMsg = "<red>Your inventory is full!</red>"
	}

	return s.WithLore(
		text.Colourf("<yellow>------------</yellow>"),
		text.Colourf("<blue>Cost:</blue> %v", resource.Name(cost)),
		"",
		text.Colourf(infoMsg),
	)
}

func editName(s item.Stack, customName string) item.Stack {
	lines := strings.Split(s.CustomName(), "\n")
	lines[0] = customName
	return s.WithCustomName(strings.Join(lines, "\n"))
}

func pickaxeTier(pl *player.Player, mode int) item.Stack {
	tiers := []struct {
		tier       item.ToolTier
		name       string
		efficiency int
		cost       int
		resource   Resource
	}{
		{item.ToolTierWood, "wooden", 1, 10, Iron},
		{item.ToolTierIron, "iron", 2, 10, Iron},
		{item.ToolTierGold, "golden", 3, 3, Gold},
		{item.ToolTierDiamond, "diamond", 3, 6, Gold},
	}
	return tieredTool(pl, mode, true, tiers)
}

func axeTier(pl *player.Player, mode int) item.Stack {
	tiers := []struct {
		tier       item.ToolTier
		name       string
		efficiency int
		cost       int
		resource   Resource
	}{
		{item.ToolTierWood, "wooden", 1, 10, Iron},
		{item.ToolTierStone, "stone", 2, 10, Iron},
		{item.ToolTierIron, "iron", 2, 3, Gold},
		{item.ToolTierDiamond, "diamond", 3, 6, Gold},
	}
	return tieredTool(pl, mode, false, tiers)
}

func tieredTool(
	pl *player.Player,
	mode int,
	isPickaxe bool,
	tiers []struct {
	tier       item.ToolTier
	name       string
	efficiency int
	cost       int
	resource   Resource
},
) item.Stack {
	maxTier := -1

	for _, stack := range pl.Inventory().Items() {
		var tier item.ToolTier
		var found bool

		if isPickaxe {
			if tool, ok := stack.Item().(item.Pickaxe); ok {
				tier = tool.Tier
				found = true
			}
		} else {
			if tool, ok := stack.Item().(item.Axe); ok {
				tier = tool.Tier
				found = true
			}
		}

		if found {
			for i, t := range tiers {
				if t.tier.Name == tier.Name && i > maxTier {
					maxTier = i
				}
			}
		}
	}

	targetTier := maxTier + mode

	if maxTier == -1 && mode < 0 {
		return item.Stack{}
	}

	if targetTier < 0 {
		targetTier = 0
	} else if targetTier >= len(tiers) {
		targetTier = len(tiers) - 1
	}

	ti := tiers[targetTier]
	owned := targetTier <= maxTier

	var stack item.Stack
	if isPickaxe {
		stack = item.NewStack(item.Pickaxe{Tier: ti.tier}, 1)
	} else {
		stack = item.NewStack(item.Axe{Tier: ti.tier}, 1)
	}

	return shopify(pl,
		stack.AsUnbreakable().
			WithEnchantments(item.NewEnchantment(enchantment.Efficiency, ti.efficiency)),
		ti.resource, ti.cost, owned, false)
}

func getCost(s item.Stack) (resource Resource, cost int) {
	if v, ok := s.Value("cost"); ok {
		cost = v.(int)
	}
	if v, ok := s.Value("resource"); ok {
		resource = Resource(v.(int))
	}

	return resource, cost
}

func canAfford(pl *player.Player, s item.Stack) bool {
	resource, cost := getCost(s)
	return pl.Inventory().ContainsItem(item.NewStack(resource.Item(), cost))
}
