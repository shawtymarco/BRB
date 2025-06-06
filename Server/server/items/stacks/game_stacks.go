package stacks

import (
	"server/server/items"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func LeaveGameStack() item.Stack {
	s := item.NewStack(item.Compass{}, 1).WithValue("action", int(items.GameSelector))
	s = s.WithCustomName(text.Colourf("<emerald>Game Selector</emerald>"))
	return s
}
