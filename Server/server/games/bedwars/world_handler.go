package bedwars

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type WorldHandler struct {
	world.NopHandler

	game *BedWars
}

func (h WorldHandler) HandleExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool) {
	*itemDropChance = 0.8
	*spawnFire = false

	filtered := (*blocks)[:0]
	for _, pos := range *blocks {
		if blocksPlaced[vec3ToString(pos.Vec3())] != nil {
			filtered = append(filtered, pos)
		}
	}
	*blocks = filtered
}
