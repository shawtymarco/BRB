package bedwars

import (
	"server/server/living"
	"server/server/utils"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type Fireball struct {
	living.NopLivingType
	*living.Living
}

func NewFireball(pos mgl64.Vec3, tx *world.Tx) *Fireball {
	t := &Fireball{}

	conf := living.Config{
		EntityType: t,
		MaxHealth:  1,
		Speed:      0.3,
		MovementComputer: &entity.MovementComputer{
			Gravity:           0.04,
			Drag:              0,
			DragBeforeGravity: false,
		},
		Drops: []living.Drop{},
	}

	s := tx.AddEntity(world.EntitySpawnOpts{Position: pos}.New(conf.EntityType, conf)).(*Fireball)
	s.SetScale(2.5, tx)

	return s
}

func (*Fireball) EncodeEntity() string { return "minecraft:fireball" }

func (*Fireball) BBox(world.Entity) cube.BBox {
	return utils.ZeroBox
}
func (f *Fireball) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	f.Living = f.NopLivingType.Open(tx, handle, data).(*living.Living)
	return f
}

func (f *Fireball) Hurt(_ float64, _ world.DamageSource) (float64, bool) {
	return 0, false
}

func (f *Fireball) Tick(tx *world.Tx, current int64) {
	if f.OnGround() {
		block.ExplosionConfig{
			Size:           5,
			ItemDropChance: 0,
			Sound:          sound.Explosion{},
			Particle:       particle.HugeExplosion{},
		}.Explode(tx, f.Position())

		utils.Panic(f.Close())
	}

	f.Living.Tick(tx, current+1)
}
