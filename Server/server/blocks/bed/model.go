package bed

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Model struct{}

func (Model) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.5625, 1)}
}

func (Model) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
