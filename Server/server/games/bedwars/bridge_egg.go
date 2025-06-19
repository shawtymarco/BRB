package bedwars

import (
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type BridgeEgg struct {
	item.Egg
	Block world.Block
}

func (e BridgeEgg) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	opts := world.EntitySpawnOpts{Position: eyePosition(user), Velocity: user.Rotation().Vec3().Mul(1.5)}
	ent := tx.World().EntityRegistry().Config().Egg(opts, user)
	tx.AddEntity(ent)
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		time.AfterFunc(3*time.Second, func() {
			ticker.Stop()
		})

		for range ticker.C {
			ent.ExecWorld(func(tx *world.Tx, ee world.Entity) {
				tx.SetBlock(cube.PosFromVec3(ee.Position()), e.Block, nil)
			})
		}
	}()

	return true
}

// eyePosition returns the position of the eyes of the entity if the entity implements entity.Eyed, or the
// actual position if it doesn't.
func eyePosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if eyed, ok := e.(interface{ EyeHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, eyed.EyeHeight()})
	}
	return pos
}
