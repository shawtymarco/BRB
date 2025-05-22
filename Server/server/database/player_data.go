package database

import (
	"time"

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
	ProtocolId string
	Statistics Statistics
	Games      Games
}

type Statistics struct {
	RankId string
	Level  int
	XP     int
	ELO    int
	Coins  int
}

func (pd PlayerData) IsRegistered() bool {
	return pd.UserId != ""
}

func (pd PlayerData) Rank() Rank {
	return RankFromName(pd.Statistics.RankId)
}

func (pd PlayerData) MaxXP() int {
	return 5000 * pd.Statistics.Level
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
	BuildFFA struct {
		Kills        int
		Deaths       int
		BlocksPlaced int
	}
	BedFight struct {
		Kills  int
		Deaths int
	}
}
