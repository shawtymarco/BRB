package bedwars

import (
	"server/server/itemutil"
	"server/server/language"
	user2 "server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.MagicMilk, MagicMilkItem{})
}

type MagicMilkItem struct {
	item.Bucket
}

func NewMagicMilkItem() item.Stack {
	return item.NewStack(MagicMilkItem{Bucket: item.Bucket{Content: item.MilkBucketContent()}}, 1).WithValue("special_item", int16(itemutil.MagicMilk))
}

func (m MagicMilkItem) Consume(_ *world.Tx, c item.Consumer) item.Stack {
	pl := c.(*player.Player)
	u := user2.GetUser(pl)
	g := Games[u.Game.ID()]
	g.trapIgnore[pl.UUID()] = true
	pl.Message(text.Colourf(language.Translate(pl).BedWars.MagicMilkEffectGive))

	time.AfterFunc(30*time.Second, func() {
		g.trapIgnore[pl.UUID()] = false
		pl.Message(text.Colourf(language.Translate(pl).BedWars.MagicMilkEffectRemove))
	})
	return item.Stack{}
}
