package item_behavior

import (
	"server/server/items"
	"server/server/ui"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
)

func init() {
	items.ItemHandlers[int(items.GameSelector)] = GameSelectorItem{}
}

type GameSelectorItem struct {
}

func (GameSelectorItem) InteractClick(_ items.ClickType, pl *player.Player, p *cube.Pos) {
	ui.NewGamesForm().SendTo(pl)
}
