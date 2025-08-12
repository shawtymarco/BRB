package listener

import (
	"server/server"
	"server/server/database"
	"server/server/itemutil"
	"server/server/itemutil/enchants"
	"server/server/language"
	"server/server/user"
	"server/server/utils"
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"

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

	if _, ok := main.Enchantment(enchants.CustomKnockBack{}); ok {
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

func HandleItemUse(ctx *player.Context) {
	pl := ctx.Val()
	u := user.GetUser(pl)
	if u.IsCooldownActive(user.Interact, 50*time.Millisecond, false, true, false) {
		return
	}

	main, off := pl.HeldItems()
	if v, ok := main.Value("special_item"); ok && main.Count() != 0 {
		if uit, ok := itemutil.SpecialItem(itemutil.Action(v.(int16))).(item.Usable); ok {
			ctx.Cancel()

			useCtx := item.UseContext{}
			uit.Use(pl.Tx(), pl, &useCtx)
			if useCtx.CountSub > 0 {
				pl.SetHeldItems(main.Grow(-1*useCtx.CountSub), off)
			}
		}
	}
}

func HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	pl := ctx.Val()
	main, off := pl.HeldItems()
	if v, ok := main.Value("special_item"); ok && main.Count() != 0 {
		ctx.Cancel()
		if uit, ok := itemutil.SpecialItem(itemutil.Action(v.(int16))).(item.UsableOnBlock); ok {
			useCtx := item.UseContext{}
			uit.UseOnBlock(pos, face, clickPos, pl.Tx(), pl, &useCtx)
			if useCtx.CountSub > 0 {
				pl.SetHeldItems(main.Grow(-1*useCtx.CountSub), off)
			}
		}
	}
}

func HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	pl := ctx.Val()
	main, off := pl.HeldItems()
	if v, ok := main.Value("special_item"); ok {
		ctx.Cancel()
		if uit, ok := itemutil.SpecialItem(itemutil.Action(v.(int16))).(ActivatedOnStartBreak); ok {
			useCtx := item.UseContext{}
			uit.OnStartBreak(pl, pos)
			if useCtx.CountSub > 0 {
				pl.SetHeldItems(main.Grow(-1*useCtx.CountSub), off)
			}
		}
	}
}

func HandleItemConsume(ctx *player.Context, s item.Stack) {
	pl := ctx.Val()
	main, off := pl.HeldItems()
	if v, ok := main.Value("special_item"); ok {
		ctx.Cancel()
		if uit, ok := itemutil.SpecialItem(itemutil.Action(v.(int16))).(item.Consumable); ok {
			pl.SetHeldItems(uit.Consume(pl.Tx(), pl), off)
		}
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

type ActivatedOnStartBreak interface {
	OnStartBreak(pl *player.Player, pos cube.Pos) bool
}
