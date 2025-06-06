package listener

import (
	"server/server"
	"server/server/items"
	"server/server/user"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	*force = server.Config.Pvp.Force
	*height = server.Config.Pvp.Height
}

func HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	*attackImmunity = time.Duration(server.Config.Pvp.HitRegistration) * time.Millisecond
}

func HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	pl := ctx.Val()
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnStartBreak, pl, &pos)
		ctx.Cancel()
	}
}

func HandlePunchAir(ctx *player.Context) {
	pl := ctx.Val()
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnPunchAir, pl, nil)
		ctx.Cancel()
	}
}

func HandleItemUse(ctx *player.Context) {
	pl := ctx.Val()
	u := user.LookupPlayer(pl)
	if u.IsCooldownActive(user.INTERACT, 50*time.Millisecond, false, false) {
		return
	}

	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnItemUse, pl, nil)
		ctx.Cancel()
	}
}
