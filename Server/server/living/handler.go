package living

import (
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/world"
)

type Handler interface {
	// HandleTick handles the entity's tick.
	HandleTick(ctx *event.Context[*Living], tx *world.Tx)
	// HandleHurt handles the entity being hurt.
	HandleHurt(ctx *event.Context[*Living], damage float64, immune bool, immunity *time.Duration, src world.DamageSource)
}

// NopHandler provides a no-op implementation of the Handler interface.
type NopHandler struct{}

var _ Handler = NopHandler{}

func (NopHandler) HandleTick(*event.Context[*Living], *world.Tx) {}

func (NopHandler) HandleHurt(*event.Context[*Living], float64, bool, *time.Duration, world.DamageSource) {
}
