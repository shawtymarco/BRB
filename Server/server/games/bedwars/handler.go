package bedwars

import (
	"fmt"
	"image/color"
	"server/server/blocks/bed"
	"server/server/database"
	"server/server/game"
	"server/server/games/lobby"
	"server/server/language"
	"server/server/listener"
	"server/server/user"
	"server/server/utils"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player/title"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/df-mc/dragonfly/server/world/sound"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var blocksPlaced = make(map[string]world.Block)

type Handler struct {
	player.NopHandler

	game *BedWars
}

func Join(pl *player.Player, tx *world.Tx, teamSize int, teamCount int, typeGame game.TypeGame, isCustom bool) {
	var bwGame *BedWars
	for _, g := range Games {
		if g.Type() == typeGame && g.Stage() == game.Waiting {
			bwGame = g
			break
		}
	}

	if bwGame == nil {
		bwGame = NewBedWars(typeGame, teamSize, teamCount, isCustom)
	}

	pl.Handle(Handler{game: bwGame})

	tx.RemoveEntity(pl)
	bwGame.World().Exec(func(tx *world.Tx) {
		tx.AddEntity(pl.H())
	})
	bwGame.AddPlayerToTeam(pl, 1)

	pl.SetGameMode(world.GameModeSurvival)
	pl.Inventory().Clear()
	pl.Armour().Clear()

	u := user.LookupPlayer(pl)
	u.Game = bwGame.Game
	switch typeGame {
	case game.TypeBedWars:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDWARS</yellow></bold>"))
		break
	case game.TypeBedFight:
		u.Scoreboard = scoreboard.New(text.Colourf("<bold><yellow>BEDFIGHT</yellow></bold>"))
		break
	default:
		panic("Unhandled game type")
	}

	pl.SetNameTag(database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Color()).Name(u.Data))
	pl.Teleport(bwGame.MapConfig().SpawnPoint)

	bwGame.ForEachPlayer(func(pl *player.Player) {
		pl.Message(text.Colourf(language.Translate(pl).Game.JoinGame, database.LobbyNameDisplay.Name(u.Data), len(bwGame.OriginalPlayers()), teamSize*teamCount))
	})
}

func (h Handler) HandleQuit(pl *player.Player) {
	user.Save(pl)
	h.game.RemovePlayerFromTeam(pl)
	lobby.Join(pl)
}

func (Handler) HandleChat(ctx *player.Context, msg *string) {
	ctx.Cancel()

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	*msg = strings.ReplaceAll(*msg, "§r", "")
	newMsg := fmt.Sprintf("%v<white>: %v</white>", database.BedWarsNameDisplay(u.Game.PlayerTeam(pl).Color()).Name(u.Data), *msg)
	*msg = text.Colourf(newMsg)

	_, _ = fmt.Fprintf(chat.Global, *msg)
}

func (Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	listener.HandleAttackEntity(ctx, e, force, height, critical)
}

func (h Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	listener.HandleHurt(ctx, damage, immune, attackImmunity, src)

	pl := ctx.Val()
	u := user.LookupPlayer(pl)

	if h.game.Stage() < game.Running {
		ctx.Cancel()
		return
	}

	if s, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := s.Attacker.(*player.Player); ok {
			ua := user.LookupPlayer(attacker)
			if pl.Health() <= *damage {
				ctx.Cancel()

				onDeath(h.game, pl, u, ua)
			}
		}
	} else if _, ok := src.(entity.VoidDamageSource); ok {
		ctx.Cancel()

		onDeath(h.game, pl, u, nil)
	} else {
		ctx.Cancel()
	}
}

func onDeath(g *BedWars, pl *player.Player, u *user.User, ua *user.User) {
	pl.Heal(pl.MaxHealth(), effect.InstantHealingSource{})
	pl.SetGameMode(world.GameModeSpectator)
	pl.Inventory().Clear()
	pl.Armour().Clear()

	finalKill := ""

	if g.PlayerTeam(pl).Status == game.BedBroken {
		finalKill = text.Colourf("<bold><aqua>FINAL KILL!</aqua></bold>")
		g.PlayerTeam(pl).RemovePlayerFromActive(pl)
	} else {
		h := pl.H()
		go func() {
			i := 5
			for range time.NewTicker(time.Second).C {
				if i == 0 {
					h.ExecWorld(func(tx *world.Tx, e world.Entity) {
						pl := e.(*player.Player)
						pl.Teleport(g.MapConfig().TeamSpawnPoints[g.PlayerTeam(pl).ID()])
						pl.SetGameMode(world.GameModeSurvival)
						giveKit(pl, g)
					})
					break
				} else {
					pl.SendTitle(title.New(text.Colourf(language.Translate(pl).BedWars.YouDiedTitle)).WithSubtitle(text.Colourf(language.Translate(pl).BedWars.YouDiedSubTitle, i)))
				}
				i--
			}
		}()
	}

	g.ForEachPlayer(func(p *player.Player) {
		if ua == nil {
			p.Message(text.Colourf(language.Translate(p).BedWars.VoidDeath, database.BedWarsNameDisplay(g.PlayerTeam(pl).Color()).Name(u.Data), finalKill))
		} else {
			p.Message(text.Colourf(language.Translate(p).BedWars.KilledBy, database.BedWarsNameDisplay(g.PlayerTeam(pl).Color()).Name(u.Data), database.BedWarsNameDisplay(g.PlayerTeam(ua.Player()).Color()).Name(ua.Data), finalKill))
		}
	})

	if ua != nil {
		ua.GameInfo.BedWars.Kills++
		ua.Player().PlaySound(sound.Experience{})
	}
}

func (h Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	pl := ctx.Val()
	if h.game.Stage() < game.Running {
		pl.Message(text.Colourf(language.Translate(pl).Game.Error.CannotBreakBlocksBecauseGameNotStarted))
		ctx.Cancel()
		return
	}

	blocksPlaced[vec3ToString(pos.Vec3())] = b
}

func (h Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	pl := ctx.Val()
	u := user.LookupPlayer(pl)
	b := pl.Tx().Block(pos)
	_, isEndstone := b.(block.EndStone)
	_, isPlank := b.(block.Planks)

	if isEndstone || isPlank {
		return
	}

	if bb, isBed := b.(bed.Bed); isBed {
		var teamIndex int
		var bedColor string
		switch bb.Colour {
		case item.ColourRed():
			teamIndex = 0
			bedColor = text.Colourf("<red>Red</red> Bed")
			break
		case item.ColourBlue():
			teamIndex = 1
			bedColor = text.Colourf("<blue>Blue</blue> Bed")
			break
		case item.ColourGreen():
			teamIndex = 2
			bedColor = text.Colourf("<green>Green</green> Bed")
			break
		case item.ColourYellow():
			teamIndex = 3
			bedColor = text.Colourf("<yellow>Yellow</yello> Bed")
			break
		}

		if h.game.PlayerTeam(pl).ID() == teamIndex {
			ctx.Cancel()
			pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakBed))
			return
		}

		h.game.Teams()[teamIndex].Status = game.BedBroken

		pl.Message(text.Colourf(language.Translate(pl).BedWars.BedBreak, bedColor, database.BedWarsNameDisplay(h.game.PlayerTeam(pl).Color()).Name(u.Data)))

		return
	}

	if blocksPlaced[vec3ToString(pos.Vec3())] == nil {
		pl.Message(text.Colourf(language.Translate(pl).BedWars.Error.CannotBreakMap))
		ctx.Cancel()
	} else {
		blocksPlaced[vec3ToString(pos.Vec3())] = nil
	}
}

func (Handler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
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

func vec3ToString(v mgl64.Vec3) string {
	return fmt.Sprintf("(%d, %d, %d)", int(v.X()), int(v.Y()), int(v.Z()))
}

func giveKit(pl *player.Player, g *BedWars) {
	teamColor := g.PlayerTeam(pl).Color()
	var c color.RGBA
	var woolColor item.Colour
	switch teamColor {
	case "red":
		c = color.RGBA{R: 255}
		woolColor = item.ColourRed()
		break
	case "blue":
		c = color.RGBA{B: 255}
		woolColor = item.ColourBlue()
		break
	case "green":
		c = color.RGBA{G: 255}
		woolColor = item.ColourGreen()
		break
	case "yellow":
		c = color.RGBA{R: 255, G: 255}
		woolColor = item.ColourYellow()
		break
	}

	utils.Panic(pl.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierWood}, 1)))
	utils.Panic(pl.Inventory().SetItem(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierWood}, 1)))
	utils.Panic(pl.Inventory().SetItem(2, item.NewStack(item.Axe{Tier: item.ToolTierWood}, 1)))
	utils.Panic(pl.Inventory().SetItem(3, item.NewStack(item.Shears{}, 1)))
	utils.Panic(pl.Inventory().SetItem(4, item.NewStack(block.Wool{Colour: woolColor}, 64)))

	pl.Armour().Set(
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: c}}, 1),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{Colour: c}}, 1),
		item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{Colour: c}}, 1),
		item.NewStack(item.Boots{Tier: item.ArmourTierLeather{Colour: c}}, 1),
	)
}
