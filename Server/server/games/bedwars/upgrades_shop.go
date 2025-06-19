package bedwars

import (
	"server/server/game"
	"server/server/living"
	"time"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server/world"
)

type UpgradesVillagerHandler struct {
	living.NopHandler

	Game   *BedWars
	Team   *game.Team
	Player *player.Player
}

func (h UpgradesVillagerHandler) HandleHurt(ctx living.Context, damage float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	ctx.Cancel()
}
