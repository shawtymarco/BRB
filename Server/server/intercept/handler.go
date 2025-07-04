package intercept

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Context = event.Context[*Conn]

var (
	handlers []Handler
)

type Handler interface {
	HandleClientPacket(ctx *Context, pk packet.Packet)
	HandleServerPacket(ctx *Context, pk packet.Packet)
}

func Hook(handler Handler) {
	handlers = append(handlers, handler)
}
