package lobby

import (
	"fmt"
	core "server/server"
	"server/server/database"
	"server/server/items/stacks"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	player.NopHandler
}

func Join(pl *player.Player) {
	u := user.LookupPlayer(pl)

	pl.Handle(Handler{})
	pl.SetNameTag(database.LobbyNameDisplay.Name(u.Data))
	pl.Teleport(core.Config.Hub.SpawnPoint)
	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()

	u.RefreshCape()

	utils.Panic(pl.Inventory().SetItem(4, stacks.GameSelectorStack()))
	utils.Panic(pl.Inventory().SetItem(0, stacks.CosmeticsStack()))

	u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BRBW</yellow></bold>"))
	u.Scoreboard.Set(0, text.Colourf("         <yellow>▷ <white>Season 1</white> ◁</yellow>"))
	u.Scoreboard.Set(1, "§0")
	u.Scoreboard.Set(2, text.Colourf("<yellow>▌</yellow> <bold><white>ELO:</white> <green>%v</green></bold>", u.Data.Statistics.ELO))
	u.Scoreboard.Set(3, text.Colourf("<yellow>▌</yellow> <bold><white>Rank:</white> <green>%v %v</green></bold>", u.Data.Statistics.ELORank().EloIcon(), u.Data.Statistics.ELORank().EloPrefix()))
	u.Scoreboard.Set(4, "§1")
	u.Scoreboard.Set(5, text.Colourf("<yellow>▌</yellow> <bold><white>Role:</white> <green>%v</green></bold>", u.Data.Rank().Prefix()))
	u.Scoreboard.Set(6, "§2")
	u.Scoreboard.Set(7, text.Colourf("<yellow>▌</yellow> <bold><white>Coins:</white> <gold>%v</gold></bold>", u.Data.Statistics.Coins))
	u.Scoreboard.Set(8, text.Colourf("<yellow>▌</yellow> <bold><white>Experience:</white> <aqua>%v</aqua></bold>", u.Data.Statistics.XP))
	u.Scoreboard.Set(9, "§3")
	u.Scoreboard.Set(10, text.Colourf("<yellow>▌</yellow> <bold><white>Total Kills:</white> <green>%v</green></bold>", u.Data.Games.TotalKills()))
	u.Scoreboard.Set(11, text.Colourf("<yellow>▌</yellow> <bold><white>Total Wins:</white> <green>%v</green></bold>", u.Data.Games.TotalKills()))
	u.Scoreboard.Set(12, "§e")
	u.Scoreboard.Set(13, text.Colourf("<yellow>ELIAGIC.CLUB</yellow>"))
	pl.SendScoreboard(u.Scoreboard)
}

func (Handler) HandleQuit(pl *player.Player) {
	user.Save(pl)
}

func (Handler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v<white>: %v</white>", database.LobbyNameDisplay.Name(u.Data), *msg)
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func (Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	ctx.Cancel()
}

func (Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	ctx.Cancel()
}

func (Handler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	listener.HandleStartBreak(ctx, pos)
}

func (Handler) HandlePunchAir(ctx *player.Context) {
	listener.HandlePunchAir(ctx)
}

func (Handler) HandleItemUse(ctx *player.Context) {
	listener.HandleItemUse(ctx)
}
