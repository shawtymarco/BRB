package lobbyitems

import (
	"server/server/itemutil"
	"server/server/ui"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"

	"github.com/df-mc/dragonfly/server/player"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.Cosmetics, CosmeticsItem{})
}

type CosmeticsItem struct {
	block.Flower
}

func (CosmeticsItem) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	ui.NewCosmeticsForm().SendTo(user.(*player.Player))
	return true
}
