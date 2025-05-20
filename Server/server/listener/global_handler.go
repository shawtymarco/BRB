package listener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server/server"
	"server/server/items"
	"server/server/user"
	"server/server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	server2 "github.com/thronesmc/matchmaking/server"
)

func PrepareMatchmakingServer(pl *player.Player) error {
	res, err := http.Get(fmt.Sprintf("http://matchmaking:8080/?identifier=%v", os.Getenv("HOSTNAME")))
	if err != nil {
		return err
	}

	defer utils.Panic(res.Body.Close())
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var MGR server.MatchmakingGetResponse
	if err := json.Unmarshal(b, &MGR); err != nil {
		return err
	}

	MGR.Data.Players = append(MGR.Data.Players, server2.Player{Username: pl.Name(), XUID: pl.UUID().String()})
	plEnc, err := json.Marshal(MGR.Data.Players)
	if err != nil {
		return err
	}

	_, err = utils.Patch(fmt.Sprintf("http://matchmaking:8080/?identifier=%v&players=%v", os.Getenv("HOSTNAME"), string(plEnc)), "application/json", nil)
	return err
}

func HandleMatchmakingServerShutdown(pl *player.Player) error {
	stand, standErr := strconv.Atoi(os.Getenv("STANDALONE"))
	if standErr == nil && stand == 0 {
		res, err := http.Get(fmt.Sprintf("http://matchmaking:8080/?identifier=%v", os.Getenv("HOSTNAME")))
		if err != nil {
			return err
		}
		defer utils.Panic(res.Body.Close())
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var MGR server.MatchmakingGetResponse
		if err := json.Unmarshal(b, &MGR); err != nil {
			return err
		}

		MGR.Data.Players = utils.Filter[server2.Player](MGR.Data.Players, func(p server2.Player) bool {
			return p.Username != pl.Name()
		})

		_, err = utils.Patch(fmt.Sprintf("http://matchmaking:8080/?identifier=%v&players=%v", os.Getenv("HOSTNAME"), MGR.Data.Players), "application/json", nil)
		if err != nil {
			return err
		}
	}

	user.Save(pl)
	defer user.Delete(pl)
	return nil
}

func HandleGlobalChat(pl *player.Player, ctx *player.Context, msg *string) {
	ctx.Cancel()
	u := user.LookupPlayer(pl)
	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v<white>:</white> <%v>%v</%v>", pl.NameTag(), u.Data.Cosmetics.ChatColor, *msg, u.Data.Cosmetics.ChatColor)
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
	if !cooldowns.CanExecute(pl, cooldowns.InteractClick, false, false, 50*time.Millisecond) {
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
	if !cooldowns.CanExecute(pl, cooldowns.InteractClick, false, false, 100*time.Millisecond) {
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
