package bed

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func init() {
	for _, dir := range cube.Directions() {
		world.RegisterBlock(Bed{Facing: dir})
	}

	for _, dir := range cube.Directions() {
		world.RegisterBlock(Bed{Facing: dir, Head: true})
	}
}

type Bed struct {
	Colour item.Colour
	Facing cube.Direction
	Head   bool
	User   item.User
}

var hashBed = block.NextHash()

func (b Bed) Hash() (uint64, uint64) {
	return hashBed | uint64(b.Facing)<<8 | uint64(boolByte(b.Head))<<10, 0
}

func (Bed) MaxCount() int {
	return 1
}

func (Bed) Model() world.BlockModel {
	return Model{}
}

func (Bed) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (b Bed) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness: 0.2,
		Harvestable: func(t item.Tool) bool {
			return true
		},
		Effective: func(t item.Tool) bool {
			return true
		},
		Drops: func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
			return []item.Stack{}
		},
		BlastResistance: 6000,
	}
}

func (b Bed) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	return false
}

func (b Bed) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, _ *item.UseContext) bool {
	return false
}

func (b Bed) EntityLand(_ cube.Pos, _ *world.World, e world.Entity, distance *float64) {
	if s, ok := e.(sneakingEntity); ok && s.Sneaking() {
		// If the entity is sneaking, the fall distance and velocity stay the same.
		return
	}
	if _, ok := e.(*player.Player); ok {
		*distance *= 0.5
	}
	if v, ok := e.(velocityEntity); ok {
		vel := v.Velocity()
		vel[1] = vel[1] * -3 / 4
		v.SetVelocity(vel)
	}
}

type sneakingEntity interface {
	Sneaking() bool
}

type velocityEntity interface {
	Velocity() mgl64.Vec3
	SetVelocity(mgl64.Vec3)
}

func (b Bed) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, _, ok := b.side(pos, tx); !ok {
		tx.SetBlock(pos, nil, nil)
	}
}

func (b Bed) EncodeItem() (name string, meta int16) {
	return "minecraft:bed", int16(b.Colour.Uint8())
}

func (b Bed) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:bed", map[string]interface{}{
		"facing_bit":   int32(horizontalDirection(b.Facing)),
		"occupied_bit": boolByte(b.User != nil),
		"head_bit":     boolByte(b.Head),
	}
}

func (b Bed) EncodeNBT() map[string]interface{} {
	return map[string]interface{}{
		"id":    "Bed",
		"color": b.Colour.Uint8(),
	}
}

func (b Bed) DecodeNBT(data map[string]interface{}) interface{} {
	b.Colour = item.Colours()[data["color"].(uint8)]
	return b
}

func (b Bed) head(pos cube.Pos, tx *world.Tx) (Bed, cube.Pos, bool) {
	headSide, headPos, ok := b.side(pos, tx)
	if !ok {
		return Bed{}, cube.Pos{}, false
	}
	if b.Head {
		headSide, headPos = b, pos
	}
	return headSide, headPos, true
}

func (b Bed) side(pos cube.Pos, tx *world.Tx) (Bed, cube.Pos, bool) {
	face := b.Facing.Face()
	if b.Head {
		face = face.Opposite()
	}

	sidePos := pos.Side(face)
	o, ok := tx.Block(sidePos).(Bed)
	return o, sidePos, ok
}

func allBeds() (beds []world.Block) {
	for _, d := range cube.Directions() {
		beds = append(beds, Bed{Facing: d})
		beds = append(beds, Bed{Facing: d, Head: true})
	}
	return
}

func boolByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func horizontalDirection(d cube.Direction) cube.Direction {
	switch d {
	case cube.South:
		return cube.North
	case cube.West:
		return cube.South
	case cube.North:
		return cube.West
	case cube.East:
		return cube.East
	}
	panic("invalid direction")
}
