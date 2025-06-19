package bedwars

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Resource int

func (r Resource) Name(cost int) string {
	switch r {
	case Iron:
		return text.Colourf("<iron>%v Iron</iron>", cost)
	case Gold:
		return text.Colourf("<gold>%v Gold</gold>", cost)
	case Diamond:
		return text.Colourf("<diamond>%v Diamond</diamond>", cost)
	case Emerald:
		return text.Colourf("<emerald>%v Emerald</emerald>", cost)
	}
	return ""
}

func (r Resource) Item() world.Item {
	switch r {
	case Iron:
		return item.IronIngot{}
	case Gold:
		return item.GoldIngot{}
	case Diamond:
		return item.Diamond{}
	case Emerald:
		return item.Emerald{}
	}
	return nil
}

func (r Resource) Block() world.Block {
	switch r {
	case Iron, Gold:
		return block.Air{}
	case Diamond:
		return block.Diamond{}
	case Emerald:
		return block.Emerald{}
	default:
		return nil
	}
}

const (
	Iron Resource = iota
	Gold
	Diamond
	Emerald
)
