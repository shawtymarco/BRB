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

// Bed is a block, allowing players to sleep to set their spawns and skip the night.
type Bed struct {
	// Colour is the colour of the bed.
	Colour item.Colour
	// Facing is the direction that the bed is facing.
	Facing cube.Direction
	// Head is true if the bed is the head side.
	Head bool
	// User is the user that is using the bed. It is only set for the Head part of the bed.
	User item.User
}

var hashBed = block.NextHash()

func (b Bed) Hash() (uint64, uint64) {
	return hashBed | uint64(b.Facing)<<8 | uint64(boolByte(b.Head))<<10, 0
}

// MaxCount always returns 1.
func (Bed) MaxCount() int {
	return 1
}

// Model ...
func (Bed) Model() world.BlockModel {
	return Model{}
}

// SideClosed ...
func (Bed) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// BreakInfo ...
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

// UseOnBlock ...
func (b Bed) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	return false
}

// Activate ...
func (b Bed) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, _ *item.UseContext) bool {
	return false
}

// EntityLand ...
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

// sneakingEntity represents an entity that can sneak.
type sneakingEntity interface {
	// Sneaking returns true if the entity is currently sneaking.
	Sneaking() bool
}

// velocityEntity represents an entity that can maintain a velocity.
type velocityEntity interface {
	// Velocity returns the current velocity of the entity.
	Velocity() mgl64.Vec3
	// SetVelocity sets the velocity of the entity.
	SetVelocity(mgl64.Vec3)
}

func (b Bed) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, _, ok := b.side(pos, tx); !ok {
		tx.SetBlock(pos, nil, nil)
	}
}

// EncodeItem ...
func (b Bed) EncodeItem() (name string, meta int16) {
	return "minecraft:bed", int16(b.Colour.Uint8())
}

// EncodeBlock ...
func (b Bed) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:bed", map[string]interface{}{
		"facing_bit":   int32(horizontalDirection(b.Facing)),
		"occupied_bit": boolByte(b.User != nil),
		"head_bit":     boolByte(b.Head),
	}
}

// EncodeNBT ...
func (b Bed) EncodeNBT() map[string]interface{} {
	return map[string]interface{}{
		"id":    "Bed",
		"color": b.Colour.Uint8(),
	}
}

// DecodeNBT ...
func (b Bed) DecodeNBT(data map[string]interface{}) interface{} {
	b.Colour = item.Colours()[data["color"].(uint8)]
	return b
}

// head returns the head side of the bed. If neither side is a head side, the third return value is false.
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

// side returns the other side of the bed. If the other side is not a bed, the third return value is false.
func (b Bed) side(pos cube.Pos, tx *world.Tx) (Bed, cube.Pos, bool) {
	face := b.Facing.Face()
	if b.Head {
		face = face.Opposite()
	}

	sidePos := pos.Side(face)
	o, ok := tx.Block(sidePos).(Bed)
	return o, sidePos, ok
}

// allBeds returns all possible beds.
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

// horizontalDirection returns the horizontal direction of the given direction. This is a legacy type still used in
// various blocks.
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
