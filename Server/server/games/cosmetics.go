package games

import (
	"server/server/items"
	"server/server/ui"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
)

func init() {
	items.ItemHandlers[int(items.Cosmetics)] = CosmeticsItem{}
}

type CosmeticsItem struct {
}

func (CosmeticsItem) InteractClick(_ items.ClickType, pl *player.Player, p *cube.Pos) {
	ui.NewCosmeticsForm().SendTo(pl)
}
