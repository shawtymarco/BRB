package bedwars

import (
	"server/server/language"
	user2 "server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type MagicMilk struct {
	item.Bucket
	game *BedWars
}

func NewMagicMilk(game *BedWars) world.Item {
	return MagicMilk{game: game, Bucket: item.Bucket{Content: item.MilkBucketContent()}}
}

func (m MagicMilk) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if m.game == nil {
		return false
	}

	pl := user.(*player.Player)
	main, off := pl.HeldItems()

	pl.SetHeldItems(main.Grow(-1), off)

	m.game.trapIgnore[pl.UUID()] = true
	user2.LookupPlayer(pl).PlaySound("random.drink", "", "", 1, 1)
	pl.Message(text.Colourf(language.Translate(pl).BedWars.MagicMilkEffectGive))

	time.AfterFunc(30*time.Second, func() {
		m.game.trapIgnore[pl.UUID()] = false
		pl.Message(text.Colourf(language.Translate(pl).BedWars.MagicMilkEffectGive))
	})
	return true
}
