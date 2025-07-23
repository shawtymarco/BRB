package bedwars

import (
	"server/server/blocks/bed"

	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type WorldHandler struct {
	world.NopHandler

	game *BedWars
}

func (h WorldHandler) HandleSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3) {
	if _, ok := s.(sound.Attack); ok {
		ctx.Cancel()
	}
}

func (h WorldHandler) HandleExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool) {
	if *itemDropChance == 1 { // Bed TNT explosion
		*itemDropChance = 0.8

		filtered := (*blocks)[:0]
		for _, pos := range *blocks {
			b := ctx.Val().Block(pos)
			if _, ok := b.(*bed.Bed); !ok && blocksPlaced[vec3ToString(pos.Vec3())] != nil {
				filtered = append(filtered, pos)
			}
		}
		*blocks = filtered
	} else { // Sudden Death Fireball/TNT Drops
		*itemDropChance = 0
	}

	*spawnFire = false
}
