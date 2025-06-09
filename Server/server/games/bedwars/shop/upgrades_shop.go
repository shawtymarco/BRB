package shop

import (
	"server/server/living"
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type UpgradesVillagerHandler struct {
	living.NopHandler
}

func (h UpgradesVillagerHandler) HandleHurt(ctx living.Context, damage float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	ctx.Cancel()
}
