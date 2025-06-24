package bedwars

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

func NewTNT(opts world.EntitySpawnOpts) *world.EntityHandle {
	conf := tntConf
	return opts.New(entity.TNTType, conf)
}

var tntConf = entity.PassiveBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
	Expire: func(e *entity.Ent, tx *world.Tx) {
		block.ExplosionConfig{
			Size:           3,
			ItemDropChance: 0.8,
			Sound:          sound.Explosion{},
			Particle:       particle.HugeExplosion{},
		}.Explode(tx, e.Position())
	},
}

type BedTNT struct {
	block.TNT
	Block world.Block
}

func (t BedTNT) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	ent := NewTNT(world.EntitySpawnOpts{Position: clickPos})
	tx.AddEntity(ent)
	return true
}
