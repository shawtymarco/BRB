package utils

import (
	"reflect"
	"unsafe"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// UpdatePrivateField sets a private field of a session to the value passed. Credits to bedrock-goephers/unsafe
func UpdatePrivateField[T any](v any, name string, value T) {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)

	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	privateFieldValue.Set(reflect.ValueOf(value))
}

// FetchPrivateField fetches a private field of a session. Credits to bedrock-goephers/unsafe
func FetchPrivateField[T any](s any, name string) T {
	reflectedValue := reflect.ValueOf(s).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	return privateFieldValue.Interface().(T)
}

//go:linkname Session github.com/df-mc/dragonfly/server/player.(*Player).session
func Session(pl *player.Player) *session.Session

//go:linkname WritePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func WritePacket(s *session.Session, pk packet.Packet)

//go:linkname EntityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func EntityRuntimeID(s *session.Session, e world.Entity) uint64

//go:linkname FirstReplaceable github.com/df-mc/dragonfly/server/block.firstReplaceable
func FirstReplaceable(tx *world.Tx, pos cube.Pos, face cube.Face, with world.Block) (cube.Pos, cube.Face, bool)

//go:linkname Place github.com/df-mc/dragonfly/server/block.place
func Place(tx *world.Tx, pos cube.Pos, b world.Block, user item.User, ctx *item.UseContext)

//go:linkname ParseEntityMetadata github.com/df-mc/dragonfly/server/session.(*Session).parseEntityMetadata
func ParseEntityMetadata(s *session.Session, e world.Entity) protocol.EntityMetadata

//go:linkname InstanceFromItem github.com/df-mc/dragonfly/server/session.instanceFromItem
func InstanceFromItem(*session.Session, item.Stack) protocol.ItemInstance
