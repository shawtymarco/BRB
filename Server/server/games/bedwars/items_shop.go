package bedwars

import (
	"server/server/game"
	"server/server/games/buildffa"
	"server/server/user"
	"server/server/utils"
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
		allShops := s.AllItemShops()
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

func (s *itemShop) itemShopBlocks() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(block.Wool{Colour: lo.If(s.team != nil, s.team.WoolColour()).Else(item.ColourWhite())}, 16), Iron, 4, false, false)
	items[20] = shopify(s.player, item.NewStack(block.StainedGlass{Colour: lo.If(s.team != nil, s.team.WoolColour()).Else(item.ColourWhite())}, 4), Iron, 8, false, false)
	items[21] = shopify(s.player, item.NewStack(block.EndStone{}, 12), Iron, 24, false, false)
	items[22] = shopify(s.player, item.NewStack(block.Ladder{}, 16), Iron, 4, false, false)
	items[23] = shopify(s.player, item.NewStack(block.Planks{Wood: lo.If(s.player != nil, user.LookupPlayer(s.player).Data.Cosmetics.SelectedWoodType).Else(block.OakWood())}, 16), Gold, 4, false, false)
	items[24] = shopify(s.player, item.NewStack(block.Obsidian{}, 4), Gold, 3, false, false)
	return items
}

func (s *itemShop) itemShopMelee() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1), Iron, 10, false, false)
	items[20] = shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1), Gold, 7, false, false)
	items[21] = shopify(s.player, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1), Emerald, 3, false, false)
	items[22] = shopify(s.player, item.NewStack(buildffa.KnockBackStick{}, 1), Gold, 5, false, false)
	return items
}

func (s *itemShop) itemShopArmour() []item.Stack {

	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1), Iron, 30, s.player.Armour().Helmet().Item().(item.Helmet).Tier == item.ArmourTierChain{}, false)
	items[20] = shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1), Gold, 12, s.player.Armour().Helmet().Item().(item.Helmet).Tier == item.ArmourTierIron{}, false)
	items[21] = shopify(s.player, item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1), Emerald, 6, s.player.Armour().Helmet().Item().(item.Helmet).Tier == item.ArmourTierDiamond{}, false)
	return items
}

func (s *itemShop) itemShopTools() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Shears{}, 1), Iron, 10, false, false)
	items[20] = pickaxeTier(s.player, s.game.pickaxeTierPlayers[s.player])
	items[21] = axeTier(s.player, s.game.axeTierPlayers[s.player])
	return items
}

func (s *itemShop) itemShopBows() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Arrow{}, 8), Gold, 2, false, false)
	items[20] = shopify(s.player, item.NewStack(item.Bow{}, 1), Gold, 12, false, false)
	items[21] = shopify(s.player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 1)), Gold, 20, false, false)
	items[22] = shopify(s.player, item.NewStack(item.Bow{}, 1).WithEnchantments(item.NewEnchantment(enchantment.Power, 2), item.NewEnchantment(enchantment.Punch, 1)), Emerald, 4, false, false)
	return items
}

func (s *itemShop) itemShopPotions() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.Potion{Type: potion.StrongLeaping()}, 1), Emerald, 1, false, false)
	items[20] = shopify(s.player, item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1), Emerald, 1, false, false)
	items[21] = shopify(s.player, item.NewStack(item.Potion{Type: potion.LongInvisibility()}, 1), Emerald, 2, false, false)
	return items
}

func (s *itemShop) itemShopUtility() []item.Stack {
	items := s.itemShopDashboard(false)
	items[19] = shopify(s.player, item.NewStack(item.GoldenApple{}, 1), Gold, 3, false, false)
	items[20] = shopify(s.player, item.NewStack(item.FireCharge{}, 1), Iron, 40, false, false)
	items[21] = shopify(s.player, item.NewStack(block.TNT{}, 1), Gold, 8, false, false)
	items[22] = shopify(s.player, item.NewStack(item.EnderPearl{}, 1), Emerald, 4, false, false)
	items[23] = shopify(s.player, item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, 1), Gold, 3, false, false)
	items[24] = editName(shopify(s.player, item.NewStack(BridgeEgg{Block: block.Wool{Colour: s.team.WoolColour()}}, 1), Emerald, 1, false, false), text.Colourf("<green>Bridge Egg</green>"))
	items[25] = shopify(s.player, item.NewStack(item.Bucket{Content: item.MilkBucketContent()}, 1), Gold, 3, false, false)
	items[28] = shopify(s.player, item.NewStack(block.Sponge{}, 4), Gold, 3, false, false)
	items[29] = shopify(s.player, item.NewStack(item.Compass{}, 1), Emerald, 2, false, false)
	return items
}

func (s *itemShop) AllItemShops() []item.Stack {
	var res []item.Stack
	for slot, i := range s.itemShopBlocks() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopMelee() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopArmour() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopTools() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopBows() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopPotions() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}
	for slot, i := range s.itemShopUtility() {
		if slot < 10 || i.Empty() {
			continue
		}
		res = append(res, i)
	}

	return res
}

func shopify(pl *player.Player, s item.Stack, resource Resource, cost int, owned bool, fullInv bool) item.Stack {
	return s.WithCustomName(text.Colourf(
		"<green>%v</green>\n<blue>Cost:</blue> %v\n\n%v",
		utils.ItemDisplay(s),
		resource.Name(cost),
		text.Colourf(lo.If(owned, "<emerald>Purchased!</emerald>").ElseIf(pl != nil && !canAfford(pl, s), "<red>You cannot afford this!</red>").ElseIf(fullInv, text.Colourf("<red>Your inventory is full!</red>")).Else("<yellow>Click to purchase!</yellow>")),
	)).WithValue("resource", int(resource)).WithValue("cost", cost)
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

func getCost(s item.Stack) (cost int, resource Resource) {
	if v, ok := s.Value("cost"); ok {
		cost = v.(int)
	}
	if v, ok := s.Value("resource"); ok {
		resource = Resource(v.(int))
	}

	return cost, resource
}

func canAfford(pl *player.Player, s item.Stack) bool {
	cost, resource := getCost(s)
	return pl.Inventory().ContainsItem(item.NewStack(resource.Item(), cost))
}
