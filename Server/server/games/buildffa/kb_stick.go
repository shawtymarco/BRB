package buildffa

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type KnockBackStick struct{}

func (s KnockBackStick) EncodeItem() (name string, meta int16) {
	return "minecraft:stick", 0
}

func (s KnockBackStick) ToolType() item.ToolType {
	return item.TypeSword
}

func (s KnockBackStick) HarvestLevel() int {
	return 0
}

func (s KnockBackStick) BaseMiningEfficiency(world.Block) float64 {
	return 0
}

type CustomKnockBack struct{}

func (CustomKnockBack) Name() string {
	return "Knockback"
}

func (CustomKnockBack) MaxLevel() int {
	return 1
}

func (CustomKnockBack) Cost(int) (int, int) {
	return 0, 1
}

func (CustomKnockBack) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

func (CustomKnockBack) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

func (CustomKnockBack) CompatibleWithItem(_ world.Item) bool {
	return true
}
