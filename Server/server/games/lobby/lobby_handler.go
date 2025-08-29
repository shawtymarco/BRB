package lobby

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	core "server/server"
	"server/server/database"
	"server/server/font"
	"server/server/games/buildffa"
	"server/server/itemutil/stacks"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/entity/effect"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item/inventory"

	"github.com/df-mc/dragonfly/server/item"

	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	player.NopHandler
}

func Join(pl *player.Player) {
	u := user.GetUser(pl)

	if b := user.ActiveBan(u.Data); b != nil {
		u.Data.Punishments.Ban(pl, b)
		return
	}

	core.Players[pl.UUID()] = pl.Name()

	if pl.Name() == "Best KoreaWW" || pl.Name() == "Studgi" {
		u.Data.Statistics.RankId = database.Owner.Shortened()
	}

	pl.Handle(Handler{})
	invHandler := listener.ChestUIHandler{Inventory: pl.Inventory(), Funcs: []func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory){
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			ctx.Cancel()
		},
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			ctx.Cancel()
		},
		func(ctx *event.Context[inventory.Holder], slot int, stack item.Stack, inv *inventory.Inventory) {
			ctx.Cancel()
		},
	}}
	pl.Inventory().Handle(invHandler)
	pl.Armour().Handle(invHandler)

	pl.SetNameTag(database.LobbyNameDisplay.Name(u.Data))
	pl.Teleport(core.Config.Hub.SpawnPoint)
	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()
	pl.Heal(20, effect.InstantHealingSource{})
	for _, e := range pl.Effects() {
		pl.RemoveEffect(e.Type())
	}

	u.RefreshCape()

	utils.Panic(pl.Inventory().SetItem(0, stacks.NewCosmeticsItem()))
	utils.Panic(pl.Inventory().SetItem(4, stacks.NewGameSelectorItem()))
	utils.Panic(pl.Inventory().SetItem(8, stacks.NewSettingsItem()))

	go func() {
		for u.Game == nil {
			inLobby := true
			<-core.MCServer.World().Exec(func(tx *world.Tx) {
				_, inLobby = u.H().Entity(tx)
			})
			if !inLobby {
				break
			}

			if !u.Data.Statistics.RankEndsIn.IsZero() && u.Data.Statistics.RankEndsIn.Before(time.Now()) {
				u.Data.Statistics.RankEndsIn = time.Time{}
				u.Data.Statistics.RankId = database.Player.Shortened()
			}

			u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BRBW</yellow></bold>"))
			u.Scoreboard.Set(1, "§0")
			u.Scoreboard.Set(2, text.Colourf("<yellow>▌</yellow> <white>ELO:</white> <green>%v</green>", u.Data.Statistics.ELO))
			u.Scoreboard.Set(3, text.Colourf("<yellow>▌</yellow> <white>Rank:</white> <green>%v %v</green>", u.Data.Statistics.ELORank().EloIcon(), u.Data.Statistics.ELORank().EloPrefix()))
			u.Scoreboard.Set(4, "§1")
			u.Scoreboard.Set(5, text.Colourf("<yellow>▌</yellow> <white>Role:</white> <grey>%v</grey>", lo.If(u.Data.Rank() == database.Player, "Player").Else(u.Data.Rank().Prefix())))
			u.Scoreboard.Set(6, "§2")
			u.Scoreboard.Set(7, text.Colourf("<yellow>▌</yellow> <white>Coins:</white> <gold>%v</gold>", u.Data.Statistics.Coins))
			u.Scoreboard.Set(8, text.Colourf("<yellow>▌</yellow> <white>Experience:</white> <aqua>%v</aqua>", u.Data.Statistics.XP))
			u.Scoreboard.Set(9, "§3")
			u.Scoreboard.Set(10, text.Colourf("<yellow>▌</yellow> <white>Total Kills:</white> <green>%v</green>", u.Data.Games.TotalKills()))
			u.Scoreboard.Set(11, text.Colourf("<yellow>▌</yellow> <white>Total Wins:</white> <green>%v</green>", u.Data.Games.TotalWins()))
			u.Scoreboard.Set(12, "§e")
			u.Scoreboard.Set(13, font.Transform(core.IP))

			u.SendScoreboard(7)

			time.Sleep(1 * time.Second)
		}
	}()
}

func (Handler) HandleQuit(pl *player.Player) {
	delete(core.Players, pl.UUID())

	pl.Inventory().Handle(nil)
	user.Save(pl)
}

func (Handler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.GetUser(pl)

	if listener.CheckChatCoolDown(pl) {
		return
	}

	msgColor := lo.If(u.Data.Rank() <= database.Booster, "white").Else("grey")

	*msg = strings.ReplaceAll(*msg, "/", "\uE000/")
	*msg = strings.ReplaceAll(*msg, "§", "")
	*msg = strings.ReplaceAll(*msg, "%", "")
	*msg = text.Colourf("%v<grey>:</grey> <%v>%v</%v>", database.LobbyNameDisplay.Name(u.Data), msgColor, *msg, msgColor)
	*msg = strings.ReplaceAll(*msg, "\uE000/", "/")

	for p := range core.MCServer.Players(pl.Tx()) {
		if um := user.GetUser(p); um.Game == nil || um.Game.ID() == buildffa.Game.ID() {
			p.Message(*msg)
		}
	}
}

func (Handler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	pl := ctx.Val()
	jumpBox := cube.Box(-105.5, 98.0, -141.5, -104.5, 102.0, -145.5)
	middle := mgl64.Vec3{-36.5, 103.0, -143.5}
	if jumpBox.Vec3Within(newPos) && !pl.Flying() {
		pl.SetVelocity(middle.Sub(newPos).Mul(0.1))
	}

	if newPos.Y() <= 0 {
		pl.Teleport(core.Config.Hub.SpawnPoint)
	}
}

func (Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	ctx.Cancel()
}

func (Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	ctx.Cancel()
}

func (Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	ctx.Cancel()
}

func (Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	ctx.Cancel()
}

func (Handler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	ctx.Cancel()
}

func (Handler) HandleItemDrop(ctx *player.Context, s item.Stack) {
	ctx.Cancel()
}

func (Handler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	listener.HandleStartBreak(ctx, pos)
}

func (Handler) HandleItemUse(ctx *player.Context) {
	listener.HandleItemUse(ctx)
}

func (Handler) HandleCommandExecution(ctx *player.Context, command cmd.Command, args []string) {
	pl := ctx.Val()
	fmt.Println(pl.Name(), " executed the command: /", command.Name(), " ", strings.Join(args, " "))
}
