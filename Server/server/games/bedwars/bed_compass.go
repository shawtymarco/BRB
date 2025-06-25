package bedwars

import (
	"server/server/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
)

type BedCompass struct {
	item.Compass
	*BedWars
}

func (c BedCompass) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pl := user.(*player.Player)
	w := tx.World()
	go func() {
		var nearestEnemy *player.Player
		c.BedWars.ForEachActivePlayer(func(p *player.Player) {
			if p.GameMode() != world.GameModeSpectator && c.BedWars.EnemyWith(p, pl) && (nearestEnemy == nil || utils.Distance(pl.Position(), p.Position()) < utils.Distance(pl.Position(), nearestEnemy.Position())) {
				nearestEnemy = p
			}
		})
		if nearestEnemy != nil {
			spawnPos := nearestEnemy.Position()
			w.SetPlayerSpawn(pl.UUID(), cube.PosFromVec3(spawnPos))
			utils.WritePacket(utils.Session(pl), &packet.SetSpawnPosition{SpawnType: packet.SpawnTypeWorld, Position: protocol.BlockPos{int32(spawnPos.X()), int32(spawnPos.Y()), int32(spawnPos.Z())}})
			pl.PlaySound(sound.Experience{})
		}
	}()

	return true
}
