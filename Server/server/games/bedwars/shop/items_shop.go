package shop

import (
	"server/server/living"
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type ItemsVillagerHandler struct {
	living.Handler
}

func (h ItemsVillagerHandler) HandleHurt(ctx living.Context, damage float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	ctx.Cancel()
}
