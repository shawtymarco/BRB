package listener

import (
	"fmt"
	"server/server/database"
	"server/server/items"
	"server/server/user"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type LobbyHandler struct {
	player.NopHandler
}

func (LobbyHandler) HandleJoin(pl *player.Player) {
	utils.Panics(user.New(pl, false))
	pl.Handle(LobbyHandler{})

	pl.SetGameMode(world.GameModeSurvival)
}

func (LobbyHandler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v<white>: %v</white>", database.LobbyNameDisplay.Name(u.Data), *msg)
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func (LobbyHandler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	pl := ctx.Val()
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnStartBreak, pl, &pos)
		ctx.Cancel()
	}
}

func (LobbyHandler) HandlePunchAir(ctx *player.Context) {
	pl := ctx.Val()
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnPunchAir, pl, nil)
		ctx.Cancel()
	}
}

func (LobbyHandler) HandleItemUse(ctx *player.Context) {
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

func (LobbyHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	pl := ctx.Val()
	u := user.LookupPlayer(pl)
	if u.IsCooldownActive(user.INTERACT, 50*time.Millisecond, false, false) {
		return
	}

	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		items.ItemHandlers[action].InteractClick(items.OnItemUseOnBlock, pl, &pos)
		ctx.Cancel()
	}
}

func (LobbyHandler) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	pl := ctx.Val()
	mainItem, _ := pl.HeldItems()
	if action, ok := mainItem.Value("action"); ok {
		action := action.(int)
		pos := cube.PosFromVec3(e.Position())
		items.ItemHandlers[action].InteractClick(items.OnItemUseOnEntity, pl, &pos)
		ctx.Cancel()
	}
}
