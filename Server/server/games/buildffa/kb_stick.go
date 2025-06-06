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
