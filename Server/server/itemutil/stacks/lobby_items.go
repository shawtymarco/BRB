package stacks

import (
	"server/server/itemutil"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func NewCosmeticsItem() item.Stack {
	s := item.NewStack(block.Flower{Type: block.Poppy()}, 1).WithValue("special_item", int16(itemutil.Cosmetics))
	s = s.WithCustomName(text.Colourf("<emerald>Cosmetics</emerald>"))
	return s
}

func NewGameSelectorItem() item.Stack {
	s := item.NewStack(item.Compass{}, 1).WithValue("special_item", int16(itemutil.GameSelector))
	s = s.WithCustomName(text.Colourf("<emerald>Game Selector</emerald>"))
	return s
}

func NewSettingsItem() item.Stack {
	s := item.NewStack(item.BlazePowder{}, 1).WithValue("special_item", int16(itemutil.Settings))
	s = s.WithCustomName(text.Colourf("<emerald>Settings</emerald>"))
	return s
}
