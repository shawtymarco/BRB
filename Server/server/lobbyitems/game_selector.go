package lobbyitems

import (
	"server/server/itemutil"
	"server/server/ui"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"

	"github.com/df-mc/dragonfly/server/player"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.GameSelector, GameSelectorItem{})
}

type GameSelectorItem struct {
	item.Compass
}

func (GameSelectorItem) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	ui.NewGamesForm().SendTo(user.(*player.Player))
	return true
}
