package bedwars

import (
	"server/server/blocks/bed"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type WorldHandler struct {
	world.NopHandler

	game *BedWars
}

func (h WorldHandler) HandleFireSpread(ctx *world.Context, from, to cube.Pos) {
	ctx.Cancel()
}

func (h WorldHandler) HandleLeavesDecay(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (h WorldHandler) HandleExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool) {
	if *itemDropChance == 1 { // Bed TNT explosion
		*itemDropChance = 0.8

		filtered := (*blocks)[:0]
		for _, pos := range *blocks {
			b := ctx.Val().Block(pos)
			_, ok1 := b.(*bed.Bed)
			_, ok2 := b.(*block.StainedGlass)
			if !ok1 && !ok2 && blocksPlaced[vec3ToString(pos.Vec3())] != nil {
				filtered = append(filtered, pos)
			}
		}
		*blocks = filtered
	} else { // Sudden Death Fireball/TNT Drops
		*itemDropChance = 0
	}

	*spawnFire = false
}
