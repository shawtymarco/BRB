package itemutil

import (
	"github.com/df-mc/dragonfly/server/world"
)

var specialItems = make(map[Action]world.Item)

func RegisterSpecialItem(action Action, it world.Item) {
	specialItems[action] = it
}

func SpecialItem(action Action) world.Item {
	return specialItems[action]
}

type Action int

const (
	GameSelector Action = iota
	Cosmetics
	Settings

	BedCompass
	BedTNT
	BridgeEgg
	MagicMilk
	SilverfishSnowball
)
