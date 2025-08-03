package listener

import (
	"server/server"
	"server/server/database"
	"server/server/items"
	"server/server/items/stacks"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"time"

	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	pl := ctx.Val()
	main, _ := pl.HeldItems()
	*force = server.Config.Pvp.Force
	*height = server.Config.Pvp.Height

	if _, ok := main.Enchantment(stacks.CustomKnockBack{}); ok {
		*force += server.Config.Pvp.Force / 2
	}
}

func HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) bool {
	if immune {
		ctx.Cancel()
		return false
	}
	*attackImmunity = time.Duration(server.Config.Pvp.HitRegistration) * time.Millisecond
	return true
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
	u := user.GetUser(pl)
	if u.IsCooldownActive(user.Interact, 50*time.Millisecond, false, true, false) {
		return
	}

	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnItemUse, pl, nil)
		ctx.Cancel()
	}
}

func CheckChatCoolDown(pl *player.Player) bool {
	u := user.GetUser(pl)
	ok1 := u.Data.Rank() > database.Partner && u.IsCooldownActive(user.Chat, 1*time.Second, false, true, true)
	var ok2 bool

	if m := user.ActiveMute(u.Data); m != nil {
		ok2 = true
		pl.Message(text.Colourf(
			language.Translate(pl).Commands.Success.Muted,
			lo.If(m.Permanent, "permanently").Else("temporarily"),
			lo.If(m.Permanent, "").Else("for "+utils.FriendlyDuration(m.EndsAt.Sub(time.Now()))),
			m.Reason,
		))
	}

	return ok1 || ok2
}
