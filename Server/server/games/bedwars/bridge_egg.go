package bedwars

import (
	"server/server/utils"
	"time"

	"github.com/df-mc/dragonfly/server/world/particle"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

func NewBridgeEgg(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := eggConf
	conf.Owner = owner.H()
	return opts.New(entity.EggType, conf)
}

var eggConf = entity.ProjectileBehaviourConfig{
	Gravity:       0.007,
	Drag:          0.001,
	Particle:      particle.EggSmash{},
	ParticleCount: 6,
}

type BridgeEgg struct {
	item.Egg
	Block world.Block
}

func (e BridgeEgg) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	ent := NewBridgeEgg(world.EntitySpawnOpts{
		Position: eyePosition(user),
		Velocity: user.Rotation().Vec3().Mul(0.5),
	}, user)

	tx.AddEntity(ent)
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		time.AfterFunc(3*time.Second, func() {
			ticker.Stop()
			ent.ExecWorld(func(tx *world.Tx, e world.Entity) {
				utils.Panic(e.Close())
			})
		})

		for range ticker.C {
			ent.ExecWorld(func(tx *world.Tx, ee world.Entity) {
				pos := cube.PosFromVec3(ee.Position())
				go func() {
					time.AfterFunc(200*time.Millisecond, func() {
						ent.ExecWorld(func(tx *world.Tx, ee world.Entity) {
							tx.SetBlock(pos, e.Block, nil)
							blocksPlaced[vec3ToString(pos.Vec3())] = e.Block
						})
					})
				}()
			})
		}
	}()

	return true
}

func eyePosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if eyed, ok := e.(interface{ EyeHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, eyed.EyeHeight()})
	}
	return pos
}
