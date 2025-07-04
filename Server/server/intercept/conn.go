package intercept

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Conn struct {
	session.Conn
	h *world.EntityHandle
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	pkt, err := c.Conn.ReadPacket()
	if err != nil {
		return pkt, err
	}

	ctx := event.C(c)
	for _, h := range handlers {
		h.HandleClientPacket(ctx, pkt)
	}

	if ctx.Cancelled() {
		return NopPacket{}, nil
	}
	return pkt, nil
}

func (c *Conn) WritePacket(pk packet.Packet) error {
	ctx := event.C(c)
	for _, h := range handlers {
		h.HandleServerPacket(ctx, pk)
	}

	if ctx.Cancelled() {
		return nil
	}
	return c.Conn.WritePacket(pk)
}

func (c *Conn) Handle() (*world.EntityHandle, bool) {
	if c.h == nil {
		if h, ok := srv.PlayerByXUID(c.Conn.IdentityData().XUID); ok {
			c.h = h
			return c.h, true
		}
		return nil, false
	}
	return c.h, true
}
