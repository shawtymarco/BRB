package bedwars

import (
	"server/server/game"
	"server/server/games/buildffa"
	"server/server/user"
	"strings"

	"github.com/df-mc/dragonfly/server/item/potion"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/item"
)

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

type itemShop struct {
	game   *BedWars
	team   *game.Team
	player *player.Player

	isQuickBuy bool
}

func (s *itemShop) itemShopDashboard(showQuickBuy bool) []item.Stack {
	items := make([]item.Stack, 54)
	items[1] = item.NewStack(block.Concrete{}, 1).WithCustomName(text.Colourf("<aqua>Blocks</aqua>"))
	items[2] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).WithCustomName(text.Colourf("<aqua>Melee</aqua>"))
	items[3] = item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Armour</aqua>"))
	items[4] = item.NewStack(item.Pickaxe{Tier: item.ToolTierStone}, 1).WithCustomName(text.Colourf("<aqua>Tools</aqua>"))
	items[5] = item.NewStack(item.Bow{}, 1).WithCustomName(text.Colourf("<aqua>Bows</aqua>"))
	items[6] = item.NewStack(block.BrewingStand{}, 1).WithCustomName(text.Colourf("<aqua>Potions</aqua>"))
	items[7] = item.NewStack(block.TNT{}, 1).WithCustomName(text.Colourf("<aqua>Utility</aqua>"))

	if showQuickBuy {
		u := user.LookupPlayer(s.player)
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
	items[53] = lo.If(s.isQuickBuy, blazePowder).Else(blazeRod)
	return items
}

func (s *itemShop) Blocks() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(block.Wool{Colour: lo.If(s.team != nil, s.team.WoolColour()).Else(item.ColourWhite())}, 16), Iron, 4, false, false)
	items[20] = shopify(s.player, item.NewStack(block.StainedGlass{Colour: lo.If(s.team != nil, s.team.WoolColour()).Else(item.ColourWhite())}, 4), Iron, 8, false, false)
	items[21] = shopify(s.player, item.NewStack(block.EndStone{}, 12), Iron, 24, false, false)
	items[22] = shopify(s.player, item.NewStack(block.Ladder{}, 16), Iron, 4, false, false)
	items[23] = shopify(s.player, item.NewStack(block.Planks{Wood: lo.If(s.player != nil, user.LookupPlayer(s.player).Data.Cosmetics.SelectedWoodType).Else(block.OakWood())}, 16), Gold, 4, false, false)
	items[24] = shopify(s.player, item.NewStack(block.Obsidian{}, 4), Gold, 3, false, false)
	return items
}

func (s *itemShop) Melee() []item.Stack {
	items := s.itemShopDashboard(false)
	i19 := shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1), Iron, 10, false, false)
	i20 := shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1), Gold, 7, false, false)
	i21 := shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1), Emerald, 3, false, false)

	if s.team.Upgrades.Sharpness > 0 {
		i19 = i19.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
		i20 = i20.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
		i21 = i21.WithEnchantments(item.NewEnchantment(enchantment.Sharpness, s.team.Upgrades.Sharpness))
	}

	items[19] = i19
	items[20] = i20
	items[21] = i21

	items[22] = shopify(s.player, item.NewStack(buildffa.KnockBackStick{}, 1), Gold, 5, false, false)
	return items
}

func (s *itemShop) Armour() []item.Stack {
	t := s.player.Armour().Boots().Item().(item.Boots).Tier
	chain := t == item.ArmourTierChain{}
	iron := t == item.ArmourTierIron{}
	diamond := t == item.ArmourTierDiamond{}
	items := s.itemShopDashboard(false)

	i19 := shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1), Iron, 30, chain || iron || diamond, false)
	i20 := shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1), Gold, 12, iron || diamond, false)
	i21 := shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1), Emerald, 6, diamond, false)

	if s.team.Upgrades.Protection != 0 {
		i19 = i19.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
		i20 = i20.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
		i21 = i21.WithEnchantments(item.NewEnchantment(enchantment.Protection, s.team.Upgrades.Protection))
	}

	items[19] = i19
	items[20] = i20
	items[21] = i21
	return items
}

func (s *itemShop) Tools() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Shears{}, 1), Iron, 10, false, false)
	items[20] = pickaxeTier(s.player, s.game.pickaxeTierPlayers[s.player])
	items[21] = axeTier(s.player, s.game.axeTierPlayers[s.player])
	return items
}

func (s *itemShop) Bows() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Arrow{}, 8), Gold, 2, false, false)
	items[20] = shopify(s.player, item.NewStack(item.Bow{}, 1), Gold, 12, false, false)
	items[21] = shopify(s.player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 1)), Gold, 20, false, false)
	items[22] = shopify(s.player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 2), item.NewEnchantment(enchantment.Punch, 1)), Emerald, 4, false, false)
	return items
}

func (s *itemShop) Potions() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Potion{Type: potion.StrongLeaping()}, 1), Emerald, 1, false, false)
	items[20] = shopify(s.player, item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1), Emerald, 1, false, false)
	items[21] = shopify(s.player, item.NewStack(item.Potion{Type: potion.LongInvisibility()}, 1), Emerald, 2, false, false)
	return items
}

func (s *itemShop) Utility() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.GoldenApple{}, 1), Gold, 3, false, false)
	items[20] = shopify(s.player, item.NewStack(SilverfishSnowball{game: s.game}, 1), Iron, 20, false, false)
	items[21] = editName(shopify(s.player, item.NewStack(block.TNT{}, 1), Gold, 8, false, false), "TNT")
	items[22] = shopify(s.player, item.NewStack(item.EnderPearl{}, 1), Emerald, 4, false, false)
	items[23] = shopify(s.player, item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, 1), Gold, 3, false, false)
	items[24] = editName(shopify(s.player, item.NewStack(BridgeEgg{Block: block.Wool{Colour: s.team.WoolColour()}}, 1), Iron, 1, false, false), text.Colourf("<green>Bridge Egg</green>")) // TODO: Change back to 1 emerald
	items[25] = shopify(s.player, item.NewStack(item.Bucket{Content: item.MilkBucketContent()}, 1), Gold, 3, false, false)
	items[28] = shopify(s.player, item.NewStack(block.Sponge{}, 4), Gold, 3, false, false)
	items[29] = shopify(s.player, item.NewStack(item.Compass{}, 1), Emerald, 2, false, false)
	return items
}

func (s *itemShop) All() []item.Stack {
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

func axeTier(pl *player.Player, tier int) item.Stack {
	switch tier {
	case 1:
		return shopify(pl, item.NewStack(item.Axe{Tier: item.ToolTierStone}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 2)), Iron, 10, false, false)
	case 2:
		return shopify(pl, item.NewStack(item.Axe{Tier: item.ToolTierIron}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 2)), Gold, 3, false, false)
	case 3:
		return shopify(pl, item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 3)), Gold, 6, false, false)
	default:
		return shopify(pl, item.NewStack(item.Axe{Tier: item.ToolTierWood}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)), Iron, 10, false, false)
	}
}

func pickaxeTier(pl *player.Player, tier int) item.Stack {
	switch tier {
	case 1:
		return shopify(pl, item.NewStack(item.Pickaxe{Tier: item.ToolTierIron}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 2)), Iron, 10, false, false)
	case 2:
		return shopify(pl, item.NewStack(item.Pickaxe{Tier: item.ToolTierGold}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 3)), Gold, 3, false, false)
	case 3:
		return shopify(pl, item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 3)), Gold, 6, false, false)
	default:
		return shopify(pl, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency, 1)), Iron, 10, false, false)
	}
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
