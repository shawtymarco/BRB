package bedwars

import (
	"server/server/itemutil"
	user2 "server/server/user"
	"server/server/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

func init() {
	itemutil.RegisterSpecialItem(itemutil.BedCompass, BedCompassItem{})
}

type BedCompassItem struct {
	item.Compass
}

func NewBedCompassItem() item.Stack {
	return item.NewStack(BedCompassItem{}, 1).WithValue("special_item", int16(itemutil.BedCompass))
}

func (c BedCompassItem) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pl := user.(*player.Player)
	u := user2.GetUser(pl)
	g := Games[u.Game.ID()]
	w := tx.World()

	var nearestEnemy *player.Player
	g.ForEachActivePlayer(func(p *player.Player) {
		if p.GameMode() != world.GameModeSpectator && g.EnemyWith(p, pl) && (nearestEnemy == nil || utils.Distance(pl.Position(), p.Position()) < utils.Distance(pl.Position(), nearestEnemy.Position())) {
			nearestEnemy = p
		}
	}, tx)

	if nearestEnemy != nil {
		spawnPos := nearestEnemy.Position()
		w.SetPlayerSpawn(pl.UUID(), cube.PosFromVec3(spawnPos))
		utils.WritePacket(utils.Session(pl), &packet.SetSpawnPosition{SpawnType: packet.SpawnTypeWorld, Position: protocol.BlockPos{int32(spawnPos.X()), int32(spawnPos.Y()), int32(spawnPos.Z())}})
		pl.PlaySound(sound.Experience{})
	}

	return true
}
