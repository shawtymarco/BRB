package blocks

import (
	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	world.RegisterItem(TripWire{})
}

type TripWire struct {
}

func (b TripWire) EncodeItem() (name string, meta int16) {
	return "minecraft:tripwire_hook", 0
}
