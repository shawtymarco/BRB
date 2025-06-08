package database

import (
	"time"

	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/google/uuid"
)

type PlayerData struct {
	Uuid   uuid.UUID
	UserId string

	AlternativeMCAccounts []string
	AlternativeDCAccounts []string

	Username   string
	Online     bool
	FirstJoin  time.Time
	LastJoin   time.Time
	DeviceOS   protocol.DeviceOS
	ProtocolId string
	Statistics Statistics
	Cosmetics  Cosmetics
	Games      Games
}

func (pd PlayerData) IsRegistered() bool {
	return pd.UserId != ""
}

func (pd PlayerData) IsTouch() bool {
	return pd.DeviceOS == protocol.DeviceAndroid || pd.DeviceOS == protocol.DeviceIOS || pd.DeviceOS == protocol.DeviceFireOS || pd.DeviceOS == protocol.DeviceWP
}

func (pd PlayerData) Rank() Rank {
	return RankFromName(pd.Statistics.RankId)
}

func (pd PlayerData) MaxXP() int {
	return 5000 * pd.Statistics.Level
}

type Statistics struct {
	RankId string
	Level  int
	XP     int
	ELO    int
	Coins  int
}

func (s Statistics) ELORank() EloRank {
	return GetEloRank(s.ELO)
}

type Cosmetics struct {
	SelectedWoodType block.WoodType
	SelectedCape     string
}

type Games struct {
	BedWars struct {
		GamesPlayed int
		MVPCount    int // TODO behavior
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
