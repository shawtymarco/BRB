package listener

import (
	"fmt"
	"server/server/items"
	"server/server/user"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func HandleGlobalChat(pl *player.Player, ctx *player.Context, msg *string) {
	ctx.Cancel()

	u := user.LookupPlayer(pl)

	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v %v<white>:</white> %v", pl.NameTag(), *msg)
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func HandleStartBreak(pl *player.Player, ctx *player.Context, pos cube.Pos) {
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnStartBreak, pl, &pos)
		ctx.Cancel()
	}
}

func HandlePunchAir(pl *player.Player, ctx *player.Context) {
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnPunchAir, pl, nil)
		ctx.Cancel()
	}
}

func HandleItemUse(pl *player.Player, ctx *player.Context) {
	u := user.LookupPlayer(pl)
	if u.IsCooldownActive(user.INTERACT, 50*time.Millisecond, false) {
		return
	}

	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnItemUse, pl, nil)
		ctx.Cancel()
	}
}

func HandleItemUseOnBlock(pl *player.Player, ctx *player.Context, pos cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	u := user.LookupPlayer(pl)
	if u.IsCooldownActive(user.INTERACT, 50*time.Millisecond, false) {
		return
	}

	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnItemUseOnBlock, pl, &pos)
		ctx.Cancel()
	}
}

func HandleItemUseOnEntity(pl *player.Player, ctx *player.Context, e world.Entity) {
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		pos := cube.PosFromVec3(e.Position())
		items.ItemHandlers[action].InteractClick(items.OnItemUseOnEntity, pl, &pos)
		ctx.Cancel()
	}
}
