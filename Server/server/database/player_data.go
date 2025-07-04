package database

import (
	"server/server/language"
	"server/server/utils"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/samber/lo"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/google/uuid"
)

type PlayerData struct {
	UUID   uuid.UUID
	UserId string

	AlternativeMCAccounts []*AltAccount

	DeviceID      string
	HashedIP      string
	IPStoredSince time.Time

	Username    string
	Online      bool
	FirstJoin   time.Time
	LastJoin    time.Time
	DeviceOS    protocol.DeviceOS
	ProtocolId  string
	Statistics  Statistics
	Cosmetics   Cosmetics
	Settings    Settings
	Punishments Punishments
	Games       Games
}

type AltAccount struct {
	UUID     uuid.UUID
	Username string
}

func (pd *PlayerData) IsRegistered() bool {
	return pd.UserId != ""
}

func (pd *PlayerData) IsTouch() bool {
	return pd.DeviceOS == protocol.DeviceAndroid || pd.DeviceOS == protocol.DeviceIOS || pd.DeviceOS == protocol.DeviceFireOS || pd.DeviceOS == protocol.DeviceWP
}

func (pd *PlayerData) Rank() Rank {
	return RankFromName(pd.Statistics.RankId)
}

func (pd *PlayerData) MaxXP() int {
	return 5000 * pd.Statistics.Level
}

func (pd *PlayerData) AddAlt(d *PlayerData) {
	for _, alt := range pd.AlternativeMCAccounts {
		if alt.UUID == d.UUID {
			return // already added
		}
	}

	pd.AlternativeMCAccounts = append(pd.AlternativeMCAccounts, &AltAccount{UUID: d.UUID, Username: d.Username})
	d.AlternativeMCAccounts = append(d.AlternativeMCAccounts, &AltAccount{UUID: pd.UUID, Username: pd.Username})
}

type Statistics struct {
	RankId     string
	RankEndsIn time.Time
	Level      int
	XP         int
	ELO        int
	Coins      int
}

func (s Statistics) ELORank() EloRank {
	return GetEloRank(s.ELO)
}

type Cosmetics struct {
	Nickname         string
	SelectedWoodType block.WoodType
	SelectedCape     string
	OwnedCapes       []string
	ELOClaimed       bool
}

type Settings struct {
	HotBarConfig   [9]HotBarCategory
	QuickBuyConfig map[int]*int
}

type HotBarCategory int

const (
	None HotBarCategory = iota
	Melee
	Blocks
	Bows
	Potions
	Utility
	Shears
	Pickaxe
	Axe
)

func (c HotBarCategory) AsStack() item.Stack {
	switch c {
	case Melee:
		return item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).WithCustomName(text.Colourf("<emerald>Melee</emerald>"))
	case Blocks:
		return item.NewStack(block.Wool{}, 1).WithCustomName(text.Colourf("<emerald>Blocks</emerald>"))
	case Bows:
		return item.NewStack(item.Bow{}, 1).WithCustomName(text.Colourf("<emerald>Bows</emerald>"))
	case Potions:
		return item.NewStack(block.BrewingStand{}, 1).WithCustomName(text.Colourf("<emerald>Potions</emerald>"))
	case Utility:
		return item.NewStack(block.TNT{}, 1).WithCustomName(text.Colourf("<emerald>Utility</emerald>"))
	case Shears:
		return item.NewStack(item.Shears{}, 1).WithCustomName(text.Colourf("<emerald>Shears</emerald>"))
	case Pickaxe:
		return item.NewStack(item.Pickaxe{Tier: item.ToolTierIron}, 1).WithCustomName(text.Colourf("<emerald>Pickaxe</emerald>"))
	case Axe:
		return item.NewStack(item.Axe{Tier: item.ToolTierIron}, 1).WithCustomName(text.Colourf("<emerald>Axe</emerald>"))
	default:
		return item.NewStack(block.StainedGlass{Colour: item.ColourRed()}, 1).WithCustomName(text.Colourf("<red>Empty Slot</red>"))
	}
}

func HotBarCategoryFromStack(stack item.Stack) HotBarCategory {
	switch {
	case stack.Equal(Melee.AsStack()):
		return Melee
	case stack.Equal(Blocks.AsStack()):
		return Blocks
	case stack.Equal(Bows.AsStack()):
		return Bows
	case stack.Equal(Potions.AsStack()):
		return Potions
	case stack.Equal(Utility.AsStack()):
		return Utility
	case stack.Equal(Shears.AsStack()):
		return Shears
	case stack.Equal(Pickaxe.AsStack()):
		return Pickaxe
	case stack.Equal(Axe.AsStack()):
		return Axe
	}
	return None
}

type Punishments struct {
	Bans  []*PunishmentData
	Mutes []*PunishmentData
}

func (p Punishments) Ban(pl *player.Player, punishment *PunishmentData) {
	pl.Disconnect(text.Colourf(
		language.Translate(pl).Commands.Success.BanDisconnect,
		lo.If(punishment.Permanent, "permanently").Else("temporarily"),
		lo.If(punishment.Permanent, "").Else(" for "+utils.FriendlyDuration(punishment.EndsAt.Sub(time.Now()))),
		punishment.Reason,
		punishment.PunishedSince.Format("Mon, Jan 2, 2006 at 3:04 PM"),
	))
}

type PunishmentData struct {
	ID            string
	PunishedBy    string
	PunishedSince time.Time
	EndsAt        time.Time
	Permanent     bool
	Reason        string

	RemovedBy string
}

type Games struct {
	BedWars struct {
		GamesPlayed int
		MVPCount    int
		Wins        int
		WinStreak   int
		Losses      int
		Kills       int
		FinalKills  int
		BedsBroken  int
		Deaths      int
	}
	BedFight struct {
		GamesPlayed int
		Wins        int
		WinStreak   int
		Losses      int
		Kills       int
		FinalKills  int
		BedsBroken  int
		Deaths      int
	}
	BuildFFA struct {
		Kills        int
		Deaths       int
		BlocksPlaced int
	}
}

func (g Games) TotalWins() int {
	return g.BedWars.Wins + g.BedFight.Wins
}

func (g Games) TotalKills() int {
	return g.BedWars.Kills + g.BedFight.Kills + g.BuildFFA.Kills
}
