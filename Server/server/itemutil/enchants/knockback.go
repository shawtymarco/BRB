package enchants

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

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
