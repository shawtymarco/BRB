package lobbyitems

import (
	"server/server/itemutil"
	"server/server/ui"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"

	"github.com/df-mc/dragonfly/server/player"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.Settings, SettingsItem{})
}

type SettingsItem struct {
	item.BlazePowder
}

func (SettingsItem) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	ui.NewSettingsForm().SendTo(user.(*player.Player))
	return true
}
