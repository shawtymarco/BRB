package listener

import (
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PacketHandler struct{}

func (h PacketHandler) HandleClientPacket(ctx *intercept.Context, pk packet.Packet) {
}

func (h PacketHandler) HandleServerPacket(_ *intercept.Context, pk packet.Packet) {}
