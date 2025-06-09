package shop

import (
	"server/server/living"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Villager struct {
	living.NopLivingType
}

func NewVillager(pos mgl64.Vec3, name string, handler living.Handler, tx *world.Tx) *living.Living {
	m := Villager{}
	conf := living.Config{
		EntityType: m,
		MaxHealth:  1,
		Speed:      0,
		Drops:      []living.Drop{},
		MovementComputer: &entity.MovementComputer{
			Gravity:           0,
			Drag:              0,
			DragBeforeGravity: false,
		},
		EyeHeight: 0.3,
		Handler:   handler,
	}
	l := tx.AddEntity(world.EntitySpawnOpts{NameTag: name, Position: pos}.New(conf.EntityType, conf)).(*living.Living)
	return l
}

func (Villager) EncodeEntity() string {
	return "minecraft:villager"
}
func (Villager) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
