package stacks

import (
	"server/server/items"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func GameSelectorStack() item.Stack {
	s := item.NewStack(item.Compass{}, 1).WithValue("action", int(items.GameSelector))
	s = s.WithCustomName(text.Colourf("<emerald>Game Selector</emerald>"))
	return s
}

func CosmeticsStack() item.Stack {
	s := item.NewStack(block.Flower{Type: block.Poppy()}, 1).WithValue("action", int(items.Cosmetics))
	s = s.WithCustomName(text.Colourf("<emerald>Cosmetics</emerald>"))
	return s
}

func SettingsStack() item.Stack {
	s := item.NewStack(item.BlazePowder{}, 1).WithValue("action", int(items.Settings))
	s = s.WithCustomName(text.Colourf("<emerald>Settings</emerald>"))
	return s
}
